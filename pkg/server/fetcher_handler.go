package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	"github.com/go-chi/chi"
)

//HandleHome is
func (s *Server) HandleHome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("home"))
}

//HandleCreateMeasure is
func (s *Server) HandleCreateMeasure(w http.ResponseWriter, r *http.Request) {
	var cm measure.CreateMeasure

	if err := json.NewDecoder(r.Body).Decode(&cm); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err := s.measures.CreateMeasure(cm)

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
	var dtos []measure.Dto

	for _, msr := range m {
		dtos = append(dtos, msr.AsDto())
	}

	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not encode measures"))
	}
	w.WriteHeader(http.StatusOK)
}

//HandleDeleteMeasure is
func (s *Server) HandleDeleteMeasure(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = s.measures.DeleteMeasure(ID)

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
	}
}
