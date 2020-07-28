package router

import (
	"net/http"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func New(s *server.Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(contentTypeJson)
	r.Use(maxBytes)
	r.Use(middleware.Logger)
	r.MethodFunc("GET", "/", s.HandleHome)
	r.MethodFunc("POST", "/api/fetcher", s.HandleCreateMeasurement)
	r.MethodFunc("GET", "/api/fetcher", s.HandleGetMeasurements)
	r.MethodFunc("GET", "/api/fetcher/id/history", s.HandleGetHisotry)
	r.MethodFunc("DELETE", "/api/fetcher/id/", s.HandleGetHisotry)
	return r
}

func contentTypeJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf8")
		next.ServeHTTP(w, r)
	})
}

//maxBytes checks body size, if its too large returns an error 413
func maxBytes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "POST payload is too large", http.StatusRequestEntityTooLarge)
			return
		}
		next.ServeHTTP(w, r)
	})
}
