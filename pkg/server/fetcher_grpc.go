package server

import (
	context "context"
	"log"
	"time"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
)

func (s *Server) GetMeasures(context.Context, *GetMeasuresRequest) (*GetMeasuresResponse, error) {
	msrs, err := s.measures.GetMeasures()
	if err != nil {
		return nil, err
	}
	var m []*Measure
	for _, msr := range msrs {
		dto := msr.AsDto()
		m = append(m, &Measure{
			ID:       int32(dto.ID),
			Interval: int32(dto.Interval),
			URL:      dto.URL,
		})
	}
	return &GetMeasuresResponse{Measures: m}, nil

}

func (s *Server) AddProbe(c context.Context, r *AddProbeRequest) (*AddProbeResponse, error) {
	m, err := s.measures.GetMeasure(int(r.MeasureID))
	if err != nil {
		return nil, err
	}
	p := measure.NewProbe(r.Response, float64(r.Duration), float32(time.Unix(r.CreatedAt, 0).Unix())/float32(time.Millisecond))
	m.AddProbe(p)
	return &AddProbeResponse{}, nil
}

func (s *Server) ListenForChanges(r *ListenForChangesRequest, stream FetcherService_ListenForChangesServer) error {
	s.streams = append(s.streams, &stream)
	for {
		select {
		case m := <-s.Add:
			s.propagate(m, Change_CREATED)
		case m := <-s.Edt:
			s.propagate(m, Change_EDITED)
		case id := <-s.Rmv:
			for _, s := range s.streams {
				err := (*s).Send(&ListenForChangesResponse{
					Change:    Change_DELETED,
					MeasureID: int32(id),
					Measure:   nil,
				})
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
}

func (s *Server) mustEmbedUnimplementedFetcherServiceServer() {

}

func (s *Server) propagate(m measure.Measure, c Change) {
	dto := m.AsDto()
	for _, s := range s.streams {
		err := (*s).Send(&ListenForChangesResponse{
			Change:    Change_CREATED,
			MeasureID: int32(dto.ID),
			Measure: &Measure{
				ID:       int32(dto.ID),
				Interval: int32(dto.Interval),
				URL:      dto.URL,
			}})
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
