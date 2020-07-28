package router

import (
	"net/http"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//New is
func New(s *server.Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(contentTypeJSON)
	r.MethodFunc("GET", "/", s.HandleHome)
	r.MethodFunc("GET", "/api/fetcher", s.HandleGetAllMeasures)
	r.MethodFunc("POST", "/api/fetcher", s.HandleCreateMeasure)
	r.MethodFunc("DELETE", "/api/fetcher/{id}", s.HandleDeleteMeasure)

	return r
}

func contentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf8")
		next.ServeHTTP(w, r)
	})
}
