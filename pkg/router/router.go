package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func New() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.MethodFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/"))
	})
	return r
}
