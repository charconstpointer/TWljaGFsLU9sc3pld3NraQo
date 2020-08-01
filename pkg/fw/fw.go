package fw

import (
	"fmt"
	"github.com/rs/zerolog/log"
)

type worker interface {
	Start() error
	AddJob()
	Done() chan struct{}
}

type Worker struct {
	jobs    []job
	queue   chan job
	d       chan struct{}
	running bool
}

func NewWorker() Worker {
	return Worker{
		jobs:  []job{},
		queue: make(chan job),
	}
}

func (w *Worker) Start() error {
	ready := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case ready <- struct{}{}:
				//close(ready)
			case job := <-w.queue:
				w.jobs = append(w.jobs, job)
				log.Info().Msg("added new job")
			}
		}
	}()
	<-ready
	w.running = true
	return nil
}

func (w *Worker) AddJob(j job) error {
	if !w.running {
		return fmt.Errorf("could not enqueue new job, please make sure that the worker is running")
	}
	select {
	case w.queue <- j:
		log.Info().Msg("enqueued new job")
		//default:
		//	return fmt.Errorf("could not enqueue new job, please make sure that the worker is running")
	}
	return nil

}

func (w *Worker) Done() chan struct{} {
	return w.d
}
