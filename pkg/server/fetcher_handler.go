package server

import (
	"encoding/json"
	"net/http"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measurement"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server/repository"
	"github.com/labstack/gommon/log"
)

func (s *Server) HandleHome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("home"))
}

//HandleCreateMeasurement handles creation of a new measurement
func (s *Server) HandleCreateMeasurement(w http.ResponseWriter, r *http.Request) {
	var m measurement.Measurement

	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	repository.CreateMeasurement(&s.measurements, &m)
	w.WriteHeader(http.StatusOK)
}

//HandleGetMeasurements returns all previously created measurements
func (s *Server) HandleGetMeasurements(w http.ResponseWriter, r *http.Request) {
	m := repository.GetMeasurements(s.measurements)
	if err := json.NewEncoder(w).Encode(&m); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

//HandleGetHistory returns history of a measurement
func (s *Server) HandleGetHisotry(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func (s *Server) HandleDeleteMeasurement(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
