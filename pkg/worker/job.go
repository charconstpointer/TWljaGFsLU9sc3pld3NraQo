package worker

import (
	"context"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"net/http"
	"time"
)

type job interface {
	Exec(ctx context.Context, res chan<- Result) error
	Stop()
	Id() int
}

type Job struct {
	id       int
	url      string
	interval int
	d        chan struct{}
}

func NewJob(id int, url string, interval int) Job {
	return Job{
		id:       id,
		url:      url,
		interval: interval,
		d:        make(chan struct{}),
	}
}

func (j Job) Exec(ctx context.Context, res chan<- Result) error {

	//i dont like this, could probably pass httpclient as a dependency, and sync.pool them in the worker?
	c := http.Client{
		Timeout: 5 * time.Second,
	}

	g := errgroup.Group{}
	g.Go(func() error {
		t := time.NewTicker(time.Duration(j.interval) * time.Second)
		for {
			select {
			case _ = <-j.d:
				//job termination
				return nil
			case _ = <-t.C:
				j.onTick(c, res)
			case _ = <-ctx.Done():
				return ctx.Err()

			}
		}
	})
	err := g.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (j Job) Stop() {
	select {
	case j.d <- struct{}{}:
		log.Info().
			Int("ID", j.id).
			Msg("stopping job")
	default:
		log.Error().
			Int("ID", j.id).
			Msg("cannot stop job")
	}
}

func (j Job) Id() int {
	return j.id
}

func (j Job) onTick(c http.Client, res chan<- Result) {
	log.Info().
		Str("URL", j.url).
		Int("interval", j.interval).
		Msg("fetching")

	start := time.Now()
	r, err := c.Get(j.url)
	stop := time.Since(start)

	if err != nil {
		log.Warn().Msgf("request to %s failed", j.url)
		return
	}

	b, _ := ioutil.ReadAll(r.Body)
	rs := Result{
		ID:      j.id,
		URL:     j.url,
		Res:     string(b),
		Dur:     stop.Seconds(),
		Success: true,
		Date:    time.Now(),
	}

	select {
	case res <- rs:
	default:
		log.Warn().
			Str("URL", j.url).
			Int("interval", j.interval).
			Msg("could not save the job's outcome, channel is full")
	}

	log.Info().
		Str("URL", j.url).
		Int("interval", j.interval).
		Msg("finished successfully")
}
