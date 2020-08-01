package fw

import (
	"context"
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
			case job := <-w.queue:
				go w.runJob(job)
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

func (w *Worker) runJob(j job) {
	result := make(chan Result, 1)
	w.jobs = append(w.jobs, j)

	go func(result <-chan Result) {
		for {
			select {
			case r := <-result:
				log.Info().
					Str("result", r.value[:25]).
					Msg("job")
			}
		}
	}(result)

	log.Info().Msg("added new job")
	err := j.Exec(context.Background(), result)

	if err != nil {
		log.Error().Msg(err.Error())
	}

}
