package fetcher

import (
	"context"
	"github.com/rs/zerolog/log"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
)

type fetcher interface {
	CreateMeasure(measure.CreateMeasure) (int, error)
	UpdateMeasure(measure *Measure, interval int) (int, error)
	DeleteMeasure() error
	GetAllMeasures() []measure.Dto
	GetHistory(measureID int) ([]measure.ProbeDto, error)
}

//Fetcher represents measures http server
type Fetcher struct {
	measures measure.Measures
	Add      chan measure.Measure
	Rmv      chan int
	Edt      chan measure.Measure
	ctx      context.Context
	streams  []*FetcherService_ListenForChangesServer
}

func (s *Fetcher) CreateMeasure(cm measure.CreateMeasure) (int, error) {
	m := measure.NewMeasure(cm.URL, cm.Interval)
	id, err := s.measures.Save(m)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	err = s.enqueue(m)
	return id, err
}

func (s *Fetcher) UpdateMeasure(measure *measure.Measure, interval int) (int, error) {
	id := measure.AsDto().ID
	err := s.measures.Update(id, interval)
	if err != nil {
		return -1, err
	}
	select {
	case s.Edt <- *measure:
	default:

	}
	return id, nil
}

func (s *Fetcher) DeleteMeasure() error {
	panic("implement me")
}

func (s *Fetcher) GetAllMeasures() []measure.Dto {
	panic("implement me")
}

func (s *Fetcher) GetHistory(measureID int) ([]measure.ProbeDto, error) {
	panic("implement me")
}

//NewFetcher creates new fetcher service
func NewFetcher(ctx context.Context, measures measure.Measures) *Fetcher {
	return &Fetcher{
		measures: measures,
		Add:      make(chan measure.Measure),
		Rmv:      make(chan int),
		Edt:      make(chan measure.Measure),
		ctx:      ctx,
	}
}
