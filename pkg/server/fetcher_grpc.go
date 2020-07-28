package server

import (
	context "context"
	"fmt"
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

func (s *Server) AddProbe(c context.Context, r *AddProbeRequest) (*AddProbeRequest, error) {
	m, err := s.measures.GetMeasure(int(r.MeasureID))
	if err != nil {
		return nil, err
	}
	p := measure.NewProbe(r.Response, float64(r.Duration), float32(time.Unix(r.CreatedAt, 0).Unix())/float32(time.Millisecond))
	m.AddProbe(p)
	return &AddProbeRequest{}, nil
}

func (s *Server) ListenForChanges(r *ListenForChangesRequest, svc FetcherService_ListenForChangesServer) error {
	for {
		select {
		case m := <-s.Add:
			dto := m.AsDto()
			err := svc.Send(&ListenForChangesResponse{
				Change:    Change_CREATED,
				MeasureID: int32(dto.ID),
				Measure: &Measure{
					ID:       int32(dto.ID),
					Interval: int32(dto.Interval),
					URL:      dto.URL,
				}})
			if err != nil {
				return err
			}
		case m := <-s.Edt:
			dto := m.AsDto()
			err := svc.Send(&ListenForChangesResponse{
				Change:    Change_EDITED,
				MeasureID: int32(dto.ID),
				Measure: &Measure{
					ID:       int32(dto.ID),
					Interval: int32(dto.Interval),
					URL:      dto.URL,
				}})
			if err != nil {
				return err
			}
		case id := <-s.Rmv:
			fmt.Println("case id := <-s.Rmv:")
			err := svc.Send(&ListenForChangesResponse{
				Change:    Change_DELETED,
				MeasureID: int32(id),
				Measure:   nil,
			})
			if err != nil {
				return err
			}
		}
	}
}

func (s *Server) mustEmbedUnimplementedFetcherServiceServer() {

}
