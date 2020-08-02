package fetcher

import (
	"context"
	"github.com/rs/zerolog/log"
	"net/http"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
)

type Fetcher interface {
	CreateOrUpdate(createMeasure measure.CreateMeasure) (int, error)
	CreateMeasure(measure.CreateMeasure) (int, error)
	UpdateMeasure(measure *measure.Measure, interval int) (int, error)
	DeleteMeasure() error
	GetAllMeasures() []measure.Dto
	GetHistory(measureID int) ([]measure.ProbeDto, error)
}

type Handler interface {
	HandleCreateMeasure(w http.ResponseWriter, r *http.Request)
	HandleGetAllMeasures(w http.ResponseWriter, _ *http.Request)
	HandleDeleteMeasure(w http.ResponseWriter, r *http.Request)
	HandleGetHistory(w http.ResponseWriter, r *http.Request)
	enqueue(m *measure.Measure) error
}

//Fetcher represents measures http server
type Fetch struct {
	measures measure.Measures
	Add      chan measure.Measure
	Rmv      chan int
	Edt      chan measure.Measure
	ctx      context.Context
	streams  []*FetcherService_ListenForChangesServer
}

func (s *Fetch) CreateMeasure(cm measure.CreateMeasure) (int, error) {
	m := measure.NewMeasure(cm.URL, cm.Interval)
	id, err := s.measures.Save(m)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	err = s.enqueue(m)
	return id, err
}

func (s *Fetch) UpdateMeasure(measure *measure.Measure, interval int) (int, error) {
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

func (s *Fetch) DeleteMeasure() error {
	panic("implement me")
}

func (s *Fetch) GetAllMeasures() []measure.Dto {
	panic("implement me")
}

func (s *Fetch) GetHistory(measureID int) ([]measure.ProbeDto, error) {
	panic("implement me")
}

//NewFetch creates new fetcher service
func NewFetch(ctx context.Context, measures measure.Measures) *Fetch {
	return &Fetch{
		measures: measures,
		Add:      make(chan measure.Measure),
		Rmv:      make(chan int),
		Edt:      make(chan measure.Measure),
		ctx:      ctx,
	}
}
