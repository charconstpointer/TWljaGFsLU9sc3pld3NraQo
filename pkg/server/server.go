package server

import (
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measurement"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	db           *sqlx.DB
	measurements []*measurement.Measurement
}

func NewServer(db *sqlx.DB) *Server {
	return &Server{measurements: make([]*measurement.Measurement, 0), db: db}
}
