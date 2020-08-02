package router

import (
	"net/http"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//New is
func New(s fetcher.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(requestSize)
	r.Use(contentTypeJSON)
	r.MethodFunc("GET", "/api/fetcher", s.HandleGetAllMeasures)
	r.MethodFunc("POST", "/api/fetcher", s.HandleCreateMeasure)
	r.MethodFunc("GET", "/api/fetcher/{id}", s.HandleGetHistory)
	r.MethodFunc("DELETE", "/api/fetcher/{id}", s.HandleDeleteMeasure)

	return r
}

func contentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf8")
		next.ServeHTTP(w, r)
	})
}

func requestSize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.ContentLength > 1000000 {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}
		next.ServeHTTP(w, r)
	})
}
