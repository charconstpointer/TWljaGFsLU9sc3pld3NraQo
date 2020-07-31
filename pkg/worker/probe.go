package worker

import (
	"context"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
	"github.com/rs/zerolog/log"
)

//Probe .
type Probe struct {
	id       int
	url      string
	interval int
}

//NewProbe .
func NewProbe(id int, url string, interval int) *Probe {
	return &Probe{
		id:       id,
		url:      url,
		interval: interval,
	}
}

//AsProbe .
func AsProbe(ms []*fetcher.Measure) []*Probe {
	probes := make([]*Probe, 0)
	for _, m := range ms {
		probes = append(probes, NewProbe(int(m.ID), m.URL, int(m.Interval)))
	}
	return probes
}

//Probes .
type Probes interface {
	All(context.Context) ([]*Probe, error)
	Add(context.Context, Result) error
	Events(context.Context) chan *fetcher.ListenForChangesResponse
}

//ProbesRepo .
type ProbesRepo struct {
	c fetcher.FetcherServiceClient
}

type ProbesRepoMock struct {
}

func (p ProbesRepoMock) All(ctx context.Context) ([]*Probe, error) {
	return nil, nil
}

func (p ProbesRepoMock) Add(ctx context.Context, result Result) error {
	return nil
}

func (p ProbesRepoMock) Events(ctx context.Context) chan *fetcher.ListenForChangesResponse {
	return make(chan *fetcher.ListenForChangesResponse)
}

//NewProbesRepo .
func NewProbesRepo(c fetcher.FetcherServiceClient) *ProbesRepo {
	return &ProbesRepo{
		c: c,
	}
}

//All fetches all currenly created jobs, this should only by used on startup, after that,
//you should be listening to incoming events and reacting to them
func (r *ProbesRepo) All(ctx context.Context) ([]*Probe, error) {
	p, err := r.c.GetMeasures(ctx, &fetcher.GetMeasuresRequest{})

	if err != nil {
		return nil, err
	}
	ps := AsProbe(p.Measures)

	return ps, nil
}

//Add persists job's result
func (r *ProbesRepo) Add(ctx context.Context, res Result) error {
	_, err := r.c.AddProbe(ctx, &fetcher.AddProbeRequest{
		MeasureID: int32(res.Probe),
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
func (r *ProbesRepo) Events(ctx context.Context) chan (*fetcher.ListenForChangesResponse) {
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
