package server

import (
	"sync"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
)

//Server represents measures http server
type Server struct {
	measures measure.Measures
	mu       sync.Mutex
}

//NewServer creates new measurement server
func NewServer(measures measure.Measures) *Server {
	return &Server{measures: measures}
}
