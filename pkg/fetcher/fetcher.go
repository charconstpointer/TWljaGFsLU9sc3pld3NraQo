package fetcher

import (
	"context"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
)

//Fetcher represents measures http server
type Fetcher struct {
	measures measure.Measures
	Add      chan measure.Measure
	Rmv      chan int
	Edt      chan measure.Measure
	ctx      context.Context
	streams  []*FetcherService_ListenForChangesServer
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
