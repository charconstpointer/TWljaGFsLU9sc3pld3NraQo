package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
)

type Probe struct {
	id       int
	url      string
	interval int
}

func NewProbe(id int, url string, interval int) *Probe {
	return &Probe{
		id:       id,
		url:      url,
		interval: interval,
	}
}

func AsProbe(ms []*fetcher.Measure) []*Probe {
	probes := make([]*Probe, 0)
	for _, m := range ms {
		probes = append(probes, NewProbe(int(m.ID), m.URL, int(m.Interval)))
	}
	fmt.Println(probes)
	return probes
}

type Probes interface {
	All() ([]*Probe, error)
	Add(res Result) error
	Events() chan (*fetcher.ListenForChangesResponse)
}

type ProbesRepo struct {
	c fetcher.FetcherServiceClient
}

func NewProbesRepo(c fetcher.FetcherServiceClient) *ProbesRepo {
	return &ProbesRepo{
		c: c,
	}
}

func (r *ProbesRepo) All() ([]*Probe, error) {
	p, err := r.c.GetMeasures(context.Background(), &fetcher.GetMeasuresRequest{})

	if err != nil {
		return nil, err
	}
	ps := AsProbe(p.Measures)

	return ps, nil
}

func (r *ProbesRepo) Add(res Result) error {
	_, err := r.c.AddProbe(context.Background(), &fetcher.AddProbeRequest{
		MeasureID: int32(res.Probe),
		CreatedAt: res.Date.Unix(),
		Duration:  float32(res.Dur),
		Response:  res.Res,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *ProbesRepo) Events() chan (*fetcher.ListenForChangesResponse) {
	ec := make(chan *fetcher.ListenForChangesResponse)
	s, err := r.c.ListenForChanges(context.Background(), &fetcher.ListenForChangesRequest{})
	if err != nil {
		log.Printf("%v", err)
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
