package worker

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
)

type worker interface {
	Start()
}

//FetcherWorker is a simpe job scheduler
type FetcherWorker struct {
	c       fetcher.FetcherServiceClient
	jobs    []*Job
	mu      sync.RWMutex
	queue   chan *fetcher.Measure
	timeout int
	Done    chan struct{}
}

//NewFetcherWorker .
func NewFetcherWorker(c fetcher.FetcherServiceClient, timeout int) *FetcherWorker {
	return &FetcherWorker{
		c:       c,
		jobs:    make([]*Job, 0),
		queue:   make(chan *fetcher.Measure),
		Done:    make(chan struct{}, 1),
		timeout: timeout,
	}
}

//Start .
func (w *FetcherWorker) Start() {

	go w.listen()
	go w.manageJobs()

	msr := w.fetchInitMsr()
	for _, m := range msr {
		w.queue <- m
	}

	<-w.Done
}

func (w *FetcherWorker) fetchInitMsr() []*fetcher.Measure {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msr, _ := w.c.GetMeasures(ctx, &fetcher.GetMeasuresRequest{})
	return msr.Measures
}

func (w *FetcherWorker) exec(ctx context.Context, j *Job) {
	measure := j.M
	log.Printf("Loaded measure : %v\n", measure)

	t := time.NewTicker(time.Duration(measure.Interval) * time.Second)
	w.mu.Lock()
	w.jobs = append(w.jobs, j)
	w.mu.Unlock()

	for {
		select {
		case _ = <-t.C:
			d, b, err := w.fetch(measure)

			if err != nil {
				log.Print(err)
				continue
			}

			err = w.saveProbe(measure.ID, d, string(b))

			if err != nil {
				log.Print(err)
				continue
			}

		case _ = <-ctx.Done():
			return
		case _ = <-j.Done:
			w.removeJob(j.M.ID)
			return
		}
	}
}

func (w *FetcherWorker) listen() {
	s, err := w.c.ListenForChanges(context.Background(), &fetcher.ListenForChangesRequest{})
	if err != nil {
		log.Printf("%v", err)
	}
	for {
		res, err := s.Recv()
		if err == io.EOF {
			log.Print("no more responses")
			return
		}
		if err != nil {
			log.Print(err)
			return
		}

		fmt.Println(res.Change)
		switch res.Change {
		case fetcher.Change_DELETED:
			for _, job := range w.jobs {
				if job.M.ID == res.MeasureID {
					job.Cancel()
				}
			}

		case fetcher.Change_CREATED:
			log.Printf("enqueuing new measure %v to be executed", res.Measure)
			w.queue <- res.Measure

		case fetcher.Change_EDITED:
			fmt.Println("EDITED")
			for _, job := range w.jobs {
				if job.M.ID == res.MeasureID {
					log.Printf("cancelling measure job %v", job.M.ID)
					job.Cancel()
					job = NewJob(res.Measure)
					w.queue <- job.M
				}
			}

		}

	}
}

func (w *FetcherWorker) manageJobs() {
	for {
		select {
		case m := <-w.queue:
			go w.exec(context.Background(), NewJob(m))
		}
	}
}

func (w *FetcherWorker) fetch(m *fetcher.Measure) (int64, []byte, error) {
	log.Printf("Fetching : %s\n", m.URL)

	var start time.Time
	var duration int64

	req, err := http.NewRequest("GET", m.URL, nil)
	c := http.Client{
		Timeout: time.Duration(w.timeout) * time.Millisecond,
	}

	start = time.Now()
	res, err := c.Do(req)
	duration = time.Since(start).Milliseconds()

	if err != nil {
		return -1, nil, err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1, nil, err
	}
	return duration, b, nil
}

func (w *FetcherWorker) removeJob(ID int32) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for i, j := range w.jobs {
		if j.M.ID == ID {
			w.jobs = append(w.jobs[:i], w.jobs[i+1:]...)
			return
		}
	}
}

func (w *FetcherWorker) saveProbe(probeID int32, duration int64, response string) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := w.c.AddProbe(ctx, &fetcher.AddProbeRequest{
		MeasureID: probeID,
		CreatedAt: time.Now().Unix(),
		Duration:  float32(float64(duration) / float64(time.Millisecond)),
		Response:  response,
	})
	if err != nil {
		return err
	}
	return nil
}
