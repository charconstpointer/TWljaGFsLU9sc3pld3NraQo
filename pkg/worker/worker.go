package worker

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type worker interface {
	AddJob(*job)
	Start(context.Context)
}
type Worker struct {
	j      []*Job
	jobs   chan (*Job)
	R      chan (*Result)
	rw     sync.RWMutex
	probes Probes
}

func NewWorker(probes Probes) *Worker {
	return &Worker{
		j:      make([]*Job, 0),
		jobs:   make(chan *Job),
		R:      make(chan *Result),
		probes: probes,
	}
}

func (w *Worker) AddJob(j *Job) {
	w.rw.Lock()
	defer w.rw.Unlock()

	w.j = append(w.j, j)

	select {
	case w.jobs <- j:
	default:
	}

}

func (w *Worker) Start(ctx context.Context) error {
	g := errgroup.Group{}
	g.Go(func() error {
		for {
			select {
			case r := <-w.R:
				fmt.Println("persisting")
				w.probes.Add(ctx, *r)
			case _ = <-ctx.Done():
				return ctx.Err()
			}
		}
	})
	g.Go(func() error {
		for {
			select {
			case j := <-w.jobs:
				go func() {
					log.Printf("starting new job %v", j)
					for {
						select {
						case _ = <-(*j).T.C:
							res, err := w.execute(j)
							if err != nil {
								fmt.Print("timeout, breaking")
								return
							}
							err = w.probes.Add(ctx, *res)
							if err != nil {
								fmt.Print("w.probes.Add")
							}
						case _ = <-(*j).D:
							log.Print("stopping worker")
							return
						case _ = <-ctx.Done():
							log.Print("stopping worker")
							return
						}
					}
				}()
			case _ = <-ctx.Done():
				return ctx.Err()
			}
		}
	})
	g.Go(func() error {
		e := w.probes.Events(ctx)
		for {
			select {
			case ev := <-e:
				p := NewProbe(int(ev.Measure.ID), ev.Measure.URL, int(ev.Measure.Interval))
				job := NewJob(p)
				w.jobs <- job
				log.Printf("enqueueing new job, %v", job)
			case _ = <-ctx.Done():
				fmt.Println("e := w.probes.Events(ctx)")
				return ctx.Err()
			}
		}

	})
	g.Go(func() error {
		p, err := w.probes.All(ctx)

		if err != nil {
			log.Println("cannot do inital fetch for already existing probes")
			return err
		}
		for _, probe := range p {
			job := NewJob(probe)
			w.jobs <- job
		}
		fmt.Println("leaving  p, err := w.probes.All(ctx)")
		return nil
	})
	err := g.Wait()
	log.Print("err := g.Wait()")
	return err
}

func (w *Worker) execute(j *Job) (*Result, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	start := time.Now()
	log.Printf("starting HTTP request %s", j.Probe().url)
	r, err := client.Get(j.p.url)
	stop := time.Since(start)
	if err != nil {
		return &Result{
			Probe:   j.p.id,
			URL:     j.p.url,
			Dur:     int(stop.Nanoseconds()),
			Success: false,
			Date:    time.Now(),
		}, nil
	}

	res, err := w.parseResp(r)
	if err != nil {
		return &Result{
			Probe:   j.p.id,
			URL:     j.p.url,
			Res:     res,
			Dur:     int(stop.Nanoseconds()),
			Success: false,
			Date:    time.Now(),
		}, nil
	}
	return &Result{
		Probe:   j.p.id,
		URL:     j.p.url,
		Res:     res,
		Dur:     int(stop.Nanoseconds()),
		Success: true,
		Date:    time.Now(),
	}, nil
}

func (w *Worker) parseResp(r *http.Response) (string, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
