package fetcher

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	"github.com/go-chi/chi"
)

//HandleHome is
func (s *Fetcher) HandleHome(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("home"))
}

//HandleCreateMeasure is
func (s *Fetcher) HandleCreateMeasure(w http.ResponseWriter, r *http.Request) {
	var cm measure.CreateMeasure

	if err := json.NewDecoder(r.Body).Decode(&cm); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	matched, err := regexp.MatchString("^http.*://", cm.URL)
	if !matched {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := measure.NewMeasure(cm.URL, cm.Interval)
	err = s.measures.CreateMeasure(m)
	go func() {
		s.Add <- *m
	}()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("something went wrong ;;"))
	}

	w.WriteHeader(http.StatusOK)
}

//HandleGetAllMeasures is
func (s *Fetcher) HandleGetAllMeasures(w http.ResponseWriter, r *http.Request) {
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
func (s *Fetcher) HandleDeleteMeasure(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = s.measures.DeleteMeasure(ID)
	select {
	case s.Rmv <- ID:
		log.Printf("notification sent  %d", ID)
	default:
		log.Println("skipping sending notification")
	}

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
	}
}

//HandleGetHistory is
func (s *Fetcher) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	m, err := s.measures.GetMeasure(ID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	var dtos []measure.ProbeDto

	for _, p := range m.Probes() {
		dtos = append(dtos, p.AsDto())
	}

	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not encode probes"))
	}
	w.WriteHeader(http.StatusOK)
}