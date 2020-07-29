package server

import (
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
)

//Server represents measures http server
type Server struct {
	measures measure.Measures
	Add      chan measure.Measure
	Rmv      chan int
	Edt      chan measure.Measure
	streams  []*FetcherService_ListenForChangesServer
}

//NewFetcher creates new fetcher service
func NewFetcher(measures measure.Measures) *Server {
	return &Server{measures: measures, Add: make(chan measure.Measure), Rmv: make(chan int), Edt: make(chan measure.Measure)}
}
