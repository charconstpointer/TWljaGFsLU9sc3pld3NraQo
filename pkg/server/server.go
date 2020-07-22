package server

import "github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measurement"

type Server struct {
	measurements []*measurement.Measurement
}

func NewServer() *Server {
	return &Server{measurements: make([]*measurement.Measurement, 0)}
}
