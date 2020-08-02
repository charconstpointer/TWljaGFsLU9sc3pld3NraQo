package worker

import (
	"context"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
	"github.com/rs/zerolog/log"
)

//Unit .
type Unit struct {
	id       int
	url      string
	interval int
}

//NewProbe .
func NewUnit(id int, url string, interval int) *Unit {
	return &Unit{
		id:       id,
		url:      url,
		interval: interval,
	}
}

//AsProbe .
func AsUnit(ms []*fetcher.Measure) []*Unit {
	probes := make([]*Unit, 0)
	for _, m := range ms {
		probes = append(probes, NewUnit(int(m.ID), m.URL, int(m.Interval)))
	}
	return probes
}

//Probes .
type Units interface {
	All(context.Context) ([]*Unit, error)
	Add(context.Context, Result) error
	Events(context.Context) chan *fetcher.ListenForChangesResponse
}

//ProbesRepo .
type UnitsRepo struct {
	c fetcher.FetcherServiceClient
}

type UnitsRepoMock struct {
}

func (p UnitsRepoMock) All(_ context.Context) ([]*Unit, error) {
	return nil, nil
}

func (p UnitsRepoMock) Add(_ context.Context, _ Result) error {
	return nil
}

func (p UnitsRepoMock) Events(_ context.Context) chan *fetcher.ListenForChangesResponse {
	return make(chan *fetcher.ListenForChangesResponse)
}

//NewProbesRepo .
func NewProbesRepo(c fetcher.FetcherServiceClient) *UnitsRepo {
	return &UnitsRepo{
		c: c,
	}
}

//All fetches all currenly created jobs, this should only by used on startup, after that,
//you should be listening to incoming events and reacting to them
func (r *UnitsRepo) All(ctx context.Context) ([]*Unit, error) {
	p, err := r.c.GetMeasures(ctx, &fetcher.GetMeasuresRequest{})

	if err != nil {
		return nil, err
	}
	ps := AsUnit(p.Measures)

	return ps, nil
}

//Add persists job's result
func (r *UnitsRepo) Add(ctx context.Context, res Result) error {
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
func (r *UnitsRepo) Events(ctx context.Context) chan *fetcher.ListenForChangesResponse {
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
