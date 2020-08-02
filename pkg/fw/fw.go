package fw

import (
	"context"
	"fmt"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type worker interface {
	Start(ctx context.Context) error
	AddJob() error
	UpdateJob(job) error
	StopJob(ID int) error
	Done() chan struct{}
	Stop() error
}

type Worker struct {
	jobs    []job
	queue   chan job
	d       chan struct{}
	running bool

	bp backplane
}

func NewWorker(bp backplane) Worker {
	return Worker{
		jobs:  []job{},
		queue: make(chan job),
		bp:    bp,
	}
}

func (w *Worker) Stop() error {
	return w.bp.Close()
}

func (w *Worker) Start(ctx context.Context) error {

	g := errgroup.Group{}

	g.Go(func() error {
		e := w.bp.Events(ctx)
		for {
			select {
			case ev := <-e:
				w.handleEvent(ev)
			case _ = <-ctx.Done():
				return ctx.Err()
			}
		}

	})

	g.Go(func() error {
		jobs, err := w.bp.Jobs(ctx)

		if err != nil {
			log.Fatal().Msgf("cannot do inital fetch for already existing units")
			return err
		}
		for _, j := range jobs {
			err := w.AddJob(j)
			if err != nil {
				err := errors.Wrap(err, "cannot enqueue job received over the wire")
				log.Warn().Msg(err.Error())
				continue
			}
		}
		return nil
	})

	g.Go(func() error {
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
	})
	err := g.Wait()

	return err
}

func (w *Worker) AddJob(j job) error {
	if !w.running {
		return fmt.Errorf("could not enqueue new job, please make sure that the worker is running")
	}
	select {
	case w.queue <- j:
		log.Info().Msg("enqueued new job")
		//default:
		//	return fmt.Errorf("could not enqueue new job")
	}
	return nil

}

func (w *Worker) UpdateJob(j job) error {
	jobID := j.Id()
	err := w.stopJob(jobID)
	if err != nil {
		return err
	}
	return w.AddJob(j)
}

func (w *Worker) StopJob(ID int) error {
	return w.stopJob(ID)
}

func (w *Worker) Done() chan struct{} {
	return w.d
}

func (w *Worker) stopJob(jobID int) error {
	for i, job := range w.jobs {
		if job.Id() == jobID {
			job.Stop()
			w.jobs = append(w.jobs[:i], w.jobs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("job with an ID %d could not be found", jobID)
}

func (w *Worker) runJob(j job) {
	result := make(chan Result, 1)
	w.jobs = append(w.jobs, j)

	go func(result <-chan Result) {
		for {
			select {
			case r := <-result:
				log.Info().
					Str("result", r.Res[:25]).
					Msg("job")
				err := w.bp.SaveResult(context.Background(), r)
				if err != nil {
					log.Error().Msgf("%s", err.Error())
				}
			}
		}
	}(result)

	log.Info().Msg("added new job")
	err := j.Exec(context.Background(), result)

	if err != nil {
		log.Error().Msg(err.Error())
	}

}

func (w *Worker) handleEvent(ev *fetcher.ListenForChangesResponse) {
	switch ev.Change {
	case fetcher.Change_CREATED:
		log.Info().Msgf("starting job %d", ev.MeasureID)
		job := NewJob(int(ev.Measure.ID), ev.Measure.URL, int(ev.Measure.Interval))
		err := w.AddJob(job)
		if err != nil {
			log.Error().Msg("cannot add new job")
		}

	case fetcher.Change_EDITED:
		log.Info().Msgf("updating job %d", ev.MeasureID)
		job := NewJob(int(ev.Measure.ID), ev.Measure.URL, int(ev.Measure.Interval))
		err := w.UpdateJob(job)
		if err != nil {
			log.Error().Msg("cannot add new job")
		}

	case fetcher.Change_DELETED:
		log.Info().Msgf("deleting job %d", ev.MeasureID)
		err := w.StopJob(int(ev.MeasureID))
		if err != nil {
			log.Error().Msg("cannot add new job")
		}
	}
}
