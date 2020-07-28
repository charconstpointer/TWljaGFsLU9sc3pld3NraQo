package server

import (
	"encoding/json"
	"net/http"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
)

//HandleHome is
func (s *Server) HandleHome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("home"))
}

//HandleCreateMeasure is
func (s *Server) HandleCreateMeasure(w http.ResponseWriter, r *http.Request) {
	err := s.measures.CreateMeasure(measure.CreateMeasure{URL: "https://foo.bar", Interval: 60})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong ;;"))
	}
	w.WriteHeader(http.StatusOK)
}

//HandleGetAllMeasures is
func (s *Server) HandleGetAllMeasures(w http.ResponseWriter, r *http.Request) {
	m, err := s.measures.GetMeasures()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong ;;"))
	}
	dtos := make([]measure.Dto, len(m))

	for _, msr := range m {
		dtos = append(dtos, msr.AsDto())
	}

	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not encode measures"))
	}
	w.WriteHeader(http.StatusOK)
}
