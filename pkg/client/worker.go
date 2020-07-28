package client

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server"
)

type worker interface {
	Listen()
}

type FetcherWorker struct {
	c    server.FetcherServiceClient
	jobs []*Job
}

func NewFetcherWorker(c server.FetcherServiceClient) *FetcherWorker {
	return &FetcherWorker{c, make([]*Job, 0)}
}

func (w *FetcherWorker) Listen() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msr, err := w.c.GetMeasures(ctx, &server.GetMeasuresRequest{})
	if err != nil {
		return err
	}
	go func(w *FetcherWorker) {
		s, err := w.c.ListenForChanges(context.Background(), &server.ListenForChangesRequest{})
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
			case server.Change_DELETED:
				for _, job := range w.jobs {
					if job.M.ID == res.MeasureID {
						job.Done <- struct{}{}
					}
				}
			}
		}

	}(w)
	for _, m := range msr.Measures {
		j := NewJob(m)
		w.jobs = append(w.jobs, j)
		go func(j *Job, fc server.FetcherServiceClient, ctx context.Context) {
			measure := j.M
			log.Printf("Loaded measure : %v\n", measure)
			t := time.NewTicker(time.Duration(measure.Interval) * time.Second)
			for {
				select {
				case _ = <-t.C:
					fetch(measure, w.c)
				case _ = <-ctx.Done():
					return
				case _ = <-j.Done:
					log.Println("Terminating")
					return
				}
			}
		}(j, w.c, context.Background())
	}

	return nil
}

func fetch(m *server.Measure, fc server.FetcherServiceClient) {
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
		log.Println(err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	fc.AddProbe(ctx, &server.AddProbeRequest{
		MeasureID: m.ID,
		CreatedAt: time.Now().Unix(),
		Duration:  float32(float64(duration) / float64(time.Millisecond)),
		Response:  string(b),
	})
}
