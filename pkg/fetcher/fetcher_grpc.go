package fetcher

import (
	"context"
	"strconv"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	"github.com/rs/zerolog/log"
)

func (s *Fetch) GetMeasures(context.Context, *GetMeasuresRequest) (*GetMeasuresResponse, error) {
	measures, err := s.measures.GetAll()
	if err != nil {
		return nil, err
	}
	var m []*Measure
	for _, msr := range measures {
		dto := msr.AsDto()
		m = append(m, &Measure{
			ID:       int32(dto.ID),
			Interval: int32(dto.Interval),
			URL:      dto.URL,
		})
	}

	return &GetMeasuresResponse{Measures: m}, nil

}

func (s *Fetch) AddProbe(_ context.Context, r *AddProbeRequest) (*AddProbeResponse, error) {
	log.Info().
		Str("measure", strconv.Itoa(int(r.MeasureID))).
		Float32("duration", r.Duration).
		Msg("received new probe result")
	p := measure.NewProbe(r.Response, r.Duration, r.CreatedAt)
	err := s.measures.SaveProbe(int(r.MeasureID), *p)
	if err != nil {
		return &AddProbeResponse{}, nil
	}
	m, err := s.measures.Get(int(r.MeasureID))

	if err != nil {
		log.Warn().Msgf("received err probe result for %d ", r.MeasureID)
		return nil, err
	}
	m.AddProbe(p)
	return &AddProbeResponse{}, nil
}
func (s *Fetch) ListenForChanges(_ *ListenForChangesRequest, stream FetcherService_ListenForChangesServer) error {
	s.streams = append(s.streams, &stream)
	for {
		select {
		case m := <-s.Add:
			s.propagate(m, Change_CREATED)
		case m := <-s.Edt:
			s.propagate(m, Change_EDITED)
		case id := <-s.Rmv:
			for _, str := range s.streams {
				err := (*str).Send(&ListenForChangesResponse{
					Change:    Change_DELETED,
					MeasureID: int32(id),
					Measure:   nil,
				})
				if err != nil {
					for i, stream := range s.streams {
						log.Warn().Msg("for i, stream := range s.streams {")
						if stream == str {
							log.Warn().Msg("if *stream == str{")
							s.streams = append(s.streams[:i], s.streams[i+1:]...)
							continue
						}
					}
					return nil
				}
			}
		case _ = <-s.ctx.Done():
			log.Info().Msg("closing grpc stream")
			return s.ctx.Err()

		}
	}
}

func (s *Fetch) mustEmbedUnimplementedImprServiceServer() {

}

func (s *Fetch) propagate(m measure.Measure, c Change) {
	dto := m.AsDto()
	for _, s := range s.streams {
		err := (*s).Send(&ListenForChangesResponse{
			Change:    c,
			MeasureID: int32(dto.ID),
			Measure: &Measure{
				ID:       int32(dto.ID),
				Interval: int32(dto.Interval),
				URL:      dto.URL,
			}})
		if err != nil {
			continue
		}
	}
}
