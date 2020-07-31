package worker

import (
	"context"
	"io/ioutil"

	"net/http"
	"sync"
	"time"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

//Worker .
type Worker struct {
	j       []*Job
	jobs    chan *Job
	R       chan *Result
	rw      sync.RWMutex
	probes  Probes
	timeout int
	clients sync.Pool
}

//NewWorker .
func NewWorker(timeout int, probes Probes) *Worker {
	return &Worker{
		j:       make([]*Job, 0),
		jobs:    make(chan *Job),
		R:       make(chan *Result),
		probes:  probes,
		timeout: timeout,
		clients: sync.Pool{
			New: func() interface{} {
				return http.Client{
					Timeout: 5 * time.Second,
				}
			},
		},
	}
}

//AddJob .
func (w *Worker) AddJob(j *Job) {
	w.rw.Lock()
	defer w.rw.Unlock()

	w.j = append(w.j, j)

	select {
	case w.jobs <- j:
	default:
	}

}

//Start .
func (w *Worker) Start(ctx context.Context) error {

	g := errgroup.Group{}

	g.Go(func() error {
		for {
			select {
			case r := <-w.R:
				err := w.probes.Add(ctx, *r)
				if err != nil {
					log.Err(err)
				}
			case _ = <-ctx.Done():
				return ctx.Err()
			}
		}
	})

	g.Go(func() error {
		for {
			select {
			case j := <-w.jobs:
				go w.initJob(ctx, j)
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
				w.handleEvent(ev)
			case _ = <-ctx.Done():
				return ctx.Err()
			}
		}

	})

	g.Go(func() error {
		p, err := w.probes.All(ctx)

		if err != nil {
			log.Fatal().Msgf("cannot do inital fetch for already existing probes")
			return err
		}
		for _, probe := range p {
			job := NewJob(probe)
			w.jobs <- job
		}
		return nil
	})

	err := g.Wait()

	return err
}

func (w *Worker) initJob(ctx context.Context, j *Job) {

	log.Info().Msgf("starting new job %d", j.p.id)

	for {
		select {
		case _ = <-(*j).T.C:
			res, err := w.exec(j)
			if err != nil {
				return
			}
			log.Info().Msgf("persisting result %s", res.URL)
			err = w.probes.Add(ctx, *res)
			if err != nil {
				log.Fatal().Msgf("could not persist %s", res.URL)
				return
			}
		case _ = <-(*j).D:
			log.Info().Msg("stopping worker")
			return
		case _ = <-ctx.Done():
			log.Info().Msg("stopping worker")
			return
		}
	}

}

func (w *Worker) exec(j *Job) (*Result, error) {
	client := w.clients.Get().(http.Client)

	start := time.Now()
	log.Info().Msgf("starting HTTP request %s", j.Probe().url)
	r, err := client.Get(j.p.url)
	stop := time.Since(start)
	w.clients.Put(client)

	if err != nil {
		log.Warn().Msgf("server failed to respond for request %s ", j.Probe().url)
		return &Result{
			Probe:   j.p.id,
			URL:     j.p.url,
			Dur:     stop.Seconds(),
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
			Dur:     stop.Seconds(),
			Success: false,
			Date:    time.Now(),
		}, nil
	}

	return &Result{
		Probe:   j.p.id,
		URL:     j.p.url,
		Res:     res,
		Dur:     stop.Seconds(),
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

func (w *Worker) findJob(ID int) (int, *Job) {
	for i, job := range w.j {
		if job.Probe().id == ID {
			return i, job
		}
	}
	return -1, nil
}

func (w *Worker) handleEvent(ev *fetcher.ListenForChangesResponse) {
	switch ev.Change {
	case fetcher.Change_CREATED:
		log.Info().Msgf("starting job %d", ev.MeasureID)
		w.enqueue(ev)
	case fetcher.Change_EDITED:
		log.Info().Msgf("updating job %d", ev.MeasureID)
		w.replace(ev)
	case fetcher.Change_DELETED:
		log.Info().Msgf("deleting job %d", ev.MeasureID)
		w.delete(ev)
	}
}

func (w *Worker) replace(ev *fetcher.ListenForChangesResponse) {
	_, j := w.findJob(int(ev.MeasureID))
	if j != nil {
		go func() {
			j.D <- struct{}{}
			p := NewProbe(
				int(ev.Measure.ID),
				ev.Measure.URL,
				int(ev.Measure.Interval),
			)
			job := NewJob(p)
			w.j = append(w.j, job)
			select {
			case w.jobs <- job:
			default:
				log.Warn().Msgf("cannot start updated job %d", ev.MeasureID)
			}
		}()
	}

}

func (w *Worker) delete(ev *fetcher.ListenForChangesResponse) {
	i, j := w.findJob(int(ev.MeasureID))
	if j != nil {
		go func() {
			j.D <- struct{}{}
			w.j = append(w.j[:i], w.j[i+1:]...)
		}()
	}
}

func (w *Worker) enqueue(ev *fetcher.ListenForChangesResponse) {
	p := NewProbe(
		int(ev.Measure.ID),
		ev.Measure.URL,
		int(ev.Measure.Interval),
	)
	job := NewJob(p)
	w.j = append(w.j, job)
	select {
	case w.jobs <- job:
	default:
		log.Warn().Msgf("cannot enqueue job %d", ev.MeasureID)
	}
}
