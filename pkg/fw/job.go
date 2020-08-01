package fw

import (
	"context"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

type job interface {
	Exec(ctx context.Context) error
	desc()
}

type Job struct {
	d Desc
}


func NewJob(d Desc) Job {
	return Job{
		d: d,
	}
}

func (j Job) desc() {
	panic("implement me")
}

func (j Job) Exec(ctx context.Context) error {
	c := http.Client{}
	g := errgroup.Group{}
	g.Go(func() error {
		t := time.NewTicker(time.Duration(j.d.interval) * time.Second)
		for {
			select {
			case _ = <-t.C:
				log.Info().
					Str("URL", j.d.url).
					Int("interval", j.d.interval).
					Msg("fetching")
				_, err := c.Get(j.d.url)
				if err != nil {
					log.Warn().Msgf("request to %s failed", j.d.url)
					continue
				}
				log.Info().
					Str("URL", j.d.url).
					Int("interval", j.d.interval).
					Msg("finished successfully")

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
