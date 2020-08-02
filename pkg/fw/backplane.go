package fw

import (
	"context"
	"fmt"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
	"github.com/rs/zerolog/log"
)

type backplane interface {
	Jobs(context.Context) ([]job, error)
	SaveResult(context.Context, Result) error
	Events(context.Context) chan *fetcher.ListenForChangesResponse
	Close() error
}

type FetcherBackplane struct {
	c fetcher.FetcherServiceClient
	d chan struct{}
}

func NewFetcherBackplane(c fetcher.FetcherServiceClient) *FetcherBackplane {
	return &FetcherBackplane{
		c: c,
		d: make(chan struct{}, 1),
	}
}

//All fetches all currenly created jobs, this should only by used on startup, after that,
//you should be listening to incoming events and reacting to them
func (r *FetcherBackplane) Jobs(ctx context.Context) ([]job, error) {
	p, err := r.c.GetMeasures(ctx, &fetcher.GetMeasuresRequest{})

	if err != nil {
		return nil, err
	}
	ps := AsUnit(p.Measures)

	return ps, nil
}

func AsUnit(ms []*fetcher.Measure) []job {
	probes := make([]job, 0)
	for _, m := range ms {
		probes = append(probes, NewJob(int(m.ID), m.URL, int(m.Interval)))
	}
	return probes
}

func (r *FetcherBackplane) Close() error {
	select {
	case r.d <- struct{}{}:
		return nil
	default:
		return fmt.Errorf("cannot close gRPC connection")
	}
}

//Add persists job's result
func (r *FetcherBackplane) SaveResult(ctx context.Context, res Result) error {
	_, err := r.c.AddProbe(ctx, &fetcher.AddProbeRequest{
		MeasureID: int32(res.ID),
		CreatedAt: float32(res.Date.Unix()),
		Duration:  float32(res.Dur),
		Response:  res.Res,
	})
	if err != nil {
		return err
	}
	return nil
}

//Events returns a channel which will contain any event published by the job producer,
//for example, new job created, job edited, job deleted, worker should react to those
//events without a need to restart
func (r *FetcherBackplane) Events(ctx context.Context) chan *fetcher.ListenForChangesResponse {
	ec := make(chan *fetcher.ListenForChangesResponse)
	s, err := r.c.ListenForChanges(ctx, &fetcher.ListenForChangesRequest{})
	if err != nil {
		log.Error().Err(err)
	}
	go func() {
		for {
			res, err := s.Recv()
			if err != nil {
				close(ec)
				return
			}
			ec <- res

		}
	}()
	return ec
}
