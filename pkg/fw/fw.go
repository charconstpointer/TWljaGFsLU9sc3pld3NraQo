package fw

import (
	"fmt"
	"golang.org/x/sync/errgroup"
)

type worker interface {
	Start() error
	AddJob()
}

type Worker struct {
	jobs  []job
	queue chan job
}

func NewWorker() Worker {
	return Worker{
		jobs:  []job{},
		queue: make(chan job),
	}
}

func (w Worker) Start() error {
	g := errgroup.Group{}

	g.Go(func() error {
		for {
			select {}
		}
	})

	err := g.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (w Worker) AddJob(j job) error {
	select {
	case w.queue <- j:
		w.jobs = append(w.jobs, j)
	default:
		return fmt.Errorf("could not enqueue new job, please make sure that the worker is running")
	}
	return nil

}
