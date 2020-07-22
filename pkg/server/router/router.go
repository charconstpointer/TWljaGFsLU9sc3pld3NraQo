package router

import (
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func New(s *server.Server) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.MethodFunc("GET", "/", s.HandleHome)
	return r
}
