package fetcher

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	"github.com/go-chi/chi"
)

//HandleCreateMeasure is
func (s Impr) HandleCreateMeasure(w http.ResponseWriter, r *http.Request) {
	var cm measure.CreateMeasure

	if err := json.NewDecoder(r.Body).Decode(&cm); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	matched, _ := regexp.MatchString("^http|https.*://", cm.URL)
	if !matched {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := s.CreateOrUpdate(cm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c := createdResponse{Id: id}
	b, err := json.Marshal(c)
	_, _ = w.Write(b)
	w.WriteHeader(http.StatusOK)

}

//HandleGetAllMeasures is
func (s Impr) HandleGetAllMeasures(w http.ResponseWriter, _ *http.Request) {
	m, err := s.GetAllMeasures()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if len(m) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("could not encode measures"))
	}

	w.WriteHeader(http.StatusOK)
}

//HandleDeleteMeasure is
func (s Impr) HandleDeleteMeasure(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = s.DeleteMeasure(ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//HandleGetHistory is
func (s Impr) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	p, err := s.GetHistory(ID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("could not encode probes"))
	}
	w.WriteHeader(http.StatusOK)
}

func (s Impr) enqueue(m *measure.Measure) error {
	go func() {
		s.Add <- *m
	}()
	return nil
}
