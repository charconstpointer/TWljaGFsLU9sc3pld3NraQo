package fetcher

import (
	context "context"
	"strconv"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	"github.com/rs/zerolog/log"
)

func (s *Fetcher) GetMeasures(context.Context, *GetMeasuresRequest) (*GetMeasuresResponse, error) {
	msrs, err := s.measures.GetAll()
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

func (s *Fetcher) AddProbe(_ context.Context, r *AddProbeRequest) (*AddProbeResponse, error) {
	log.Info().
		Str("measure", strconv.Itoa(int(r.MeasureID))).
		Float32("duration", r.Duration).
		Msg("received new probe result")

	m, err := s.measures.Get(int(r.MeasureID))
	if err != nil {
		log.Warn().Msgf("received err probe result for %d ", r.MeasureID)
		return nil, err
	}
	p := measure.NewProbe(r.Response, r.Duration, r.CreatedAt)
	m.AddProbe(p)
	return &AddProbeResponse{}, nil
}
func (s *Fetcher) ListenForChanges(str FetcherService_ListenForChangesServer) error {
	go func() {
		for {
			_, err := str.Recv()
			if err != nil {
				log.Error().Err(err)
				for i, stream := range s.streams {
					log.Warn().Msg("for i, stream := range s.streams {")
					if *stream == str {
						log.Warn().Msg("if *stream == str{")
						s.streams = append(s.streams[:i], s.streams[i+1:]...)
						continue
					}
				}
				return
			}

		}
	}()
	s.streams = append(s.streams, &str)
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
					log.Fatal().Msg(err.Error())
					continue
				}
			}
		case _ = <-s.ctx.Done():
			log.Info().Msg("closing grpc stream")
			return s.ctx.Err()

		}
	}
}

// func (s *Fetcher) ListenForChanges(_ *ListenForChangesRequest, stream FetcherService_ListenForChangesServer) error {
// 	s.streams = append(s.streams, &stream)
// 	for {
// 		select {
// 		case m := <-s.Add:
// 			s.propagate(m, Change_CREATED)
// 		case m := <-s.Edt:
// 			s.propagate(m, Change_EDITED)
// 		case id := <-s.Rmv:
// 			for _, s := range s.streams {
// 				err := (*s).Send(&ListenForChangesResponse{
// 					Change:    Change_DELETED,
// 					MeasureID: int32(id),
// 					Measure:   nil,
// 				})
// 				if err != nil {
// 					log.Fatal().Msg(err.Error())
// 					continue
// 				}
// 			}
// 		case _ = <-s.ctx.Done():
// 			log.Info().Msg("closing grpc stream")
// 			return s.ctx.Err()

// 		}
// 	}
// }

func (s *Fetcher) mustEmbedUnimplementedFetcherServiceServer() {

}

func (s *Fetcher) propagate(m measure.Measure, c Change) {
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
