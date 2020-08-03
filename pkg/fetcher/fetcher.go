package fetcher

import (
	"context"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Fetch struct {
	measures measure.Measures
	Add      chan measure.Measure
	Rmv      chan int
	Edt      chan measure.Measure
	ctx      context.Context
	streams  []*FetcherService_ListenForChangesServer
}

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

func (s Fetch) mustEmbedUnimplementedFetcherServiceServer() {
	panic("implement me")
}

func NewImpr(ctx context.Context, measures measure.Measures) *Fetch {
	return &Fetch{
		measures: measures,
		Add:      make(chan measure.Measure),
		Rmv:      make(chan int),
		Edt:      make(chan measure.Measure),
		ctx:      ctx,
	}
}

func (s Fetch) CreateOrUpdate(createMeasure measure.CreateMeasure) (int, error) {
	m, _ := s.measures.GetByUrl(createMeasure.URL)
	if m != nil {
		mID := m.AsDto().ID

		err := s.measures.Update(mID, createMeasure.Interval)
		if err != nil {
			return -1, err
		}
		select {
		case s.Edt <- *m:
			log.Info().Msgf("notification sent  %d", m.ID)
		default:
			log.Info().Msgf("skipping sending notification")
		}
		return mID, nil
	}
	m = measure.NewMeasure(createMeasure.URL, createMeasure.Interval)
	ID, err := s.measures.Save(m)
	if err != nil {
		return ID, err
	}

	m.ID = ID
	select {
	case s.Add <- *m:
		log.Info().Msgf("notification sent  %d", ID)
	default:
		log.Info().Msgf("skipping sending notification")
	}
	return m.ID, nil
}

func (s Fetch) CreateMeasure(m measure.CreateMeasure) (int, error) {
	nm := measure.NewMeasure(m.URL, m.Interval)
	id, err := s.measures.Save(nm)
	return id, err
}

func (s Fetch) UpdateMeasure(measure *measure.Measure, interval int) (int, error) {
	return s.UpdateMeasure(measure, interval)
}

func (s Fetch) DeleteMeasure(ID int) error {
	err := s.measures.Delete(ID)
	if err != nil {
		return err
	}
	select {

	case s.Rmv <- ID:
		log.Info().Msgf("notification sent  %d", ID)
	default:
		log.Info().Msgf("skipping sending notification")
	}
	return nil
}

func (s Fetch) GetAllMeasures() ([]measure.Dto, error) {
	m, err := s.measures.GetAll()
	if err != nil {
		return nil, err
	}
	var dtos []measure.Dto

	for _, msr := range m {
		dtos = append(dtos, msr.AsDto())
	}
	return dtos, nil
}

func (s Fetch) GetHistory(measureID int) ([]measure.ProbeDto, error) {
	m, err := s.measures.Get(measureID)
	if err != nil {
		return nil, err
	}

	var dtos []measure.ProbeDto
	for _, p := range m.Probes {
		dtos = append(dtos, p.AsDto())
	}
	return dtos, nil
}
