package fetcher

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

//HandleCreateMeasure is
func (s *Fetcher) HandleCreateMeasure(w http.ResponseWriter, r *http.Request) {
	var cm measure.CreateMeasure

	if err := json.NewDecoder(r.Body).Decode(&cm); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	matched, _ := regexp.MatchString("^http|https.*://", cm.URL)
	if !matched {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m, _ := s.measures.GetByUrl(cm.URL)
	//if already exists then update
	if m != nil {
		id, err := s.UpdateMeasure(m, cm.Interval)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		c := createdResponse{Id: id}
		b, err := json.Marshal(c)
		_, _ = w.Write(b)
		w.WriteHeader(http.StatusOK)
		return
	}

	id, err := s.CreateMeasure(cm)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	c := createdResponse{Id: id}
	b, err := json.Marshal(c)
	_, _ = w.Write(b)
	w.WriteHeader(http.StatusOK)
}

//HandleGetAllMeasures is
func (s *Fetcher) HandleGetAllMeasures(w http.ResponseWriter, _ *http.Request) {
	m, err := s.measures.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if len(m) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	var dtos []measure.Dto

	for _, msr := range m {
		dtos = append(dtos, msr.AsDto())
	}

	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("could not encode measures"))
	}
	w.WriteHeader(http.StatusOK)
}

//HandleDeleteMeasure is
func (s *Fetcher) HandleDeleteMeasure(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = s.measures.Delete(ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	select {
	case s.Rmv <- ID:
		log.Info().Msgf("notification sent  %d", ID)
	default:
		log.Info().Msgf("skipping sending notification")
	}

	w.WriteHeader(http.StatusOK)

}

//HandleGetHistory is
func (s *Fetcher) HandleGetHistory(w http.ResponseWriter, r *http.Request) {

	ID, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	m, err := s.measures.Get(ID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var dtos []measure.ProbeDto

	for _, p := range m.Probes() {
		dtos = append(dtos, p.AsDto())
	}

	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("could not encode probes"))
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Fetcher) enqueue(m *measure.Measure) error {
	go func() {
		s.Add <- *m
	}()
	return nil
}

type createdResponse struct {
	Id int `json:"id"`
}
