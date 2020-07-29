package worker

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
)

type worker interface {
	Listen()
}

//FetcherWorker is a simpe job scheduler
type FetcherWorker struct {
	c     fetcher.FetcherServiceClient
	jobs  []*Job
	queue chan *fetcher.Measure
	Done  chan struct{}
}

//NewFetcherWorker .
func NewFetcherWorker(c fetcher.FetcherServiceClient) *FetcherWorker {
	return &FetcherWorker{c, make([]*Job, 0), make(chan *fetcher.Measure), make(chan struct{}, 1)}
}

//Start .
func (w *FetcherWorker) Start() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msr, _ := w.c.GetMeasures(ctx, &fetcher.GetMeasuresRequest{})
	// if err != nil {
	// 	return err
	// }
	go w.listen()
	go w.manageJobs()
	for _, m := range msr.Measures {
		w.queue <- m
	}

	<-w.Done
}

func (w *FetcherWorker) exec(ctx context.Context, j *Job) {
	measure := j.M
	log.Printf("Loaded measure : %v\n", measure)
	t := time.NewTicker(time.Duration(measure.Interval) * time.Second)
	w.jobs = append(w.jobs, j)
	for {
		select {
		case _ = <-t.C:
			d, b, err := fetch(measure, w.c)
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
			log.Printf("terminating measure job %v", j.M)
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
			for _, job := range w.jobs {
				if job.M.ID == res.MeasureID {
					log.Printf("cancelling measure job %v", res.Measure)
					job.Cancel()
					job = NewJob(res.Measure)
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

func fetch(m *fetcher.Measure, fc fetcher.FetcherServiceClient) (int64, []byte, error) {
	log.Printf("Fetching : %s\n", m.URL)

	var start time.Time
	var duration int64

	req, err := http.NewRequest("GET", m.URL, nil)
	c := http.Client{
		Timeout: time.Duration(m.Interval) * time.Second,
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
