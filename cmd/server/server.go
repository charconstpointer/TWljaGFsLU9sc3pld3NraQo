package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on for incoming http requests")
	flag.Parse()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	addr := ":" + strconv.Itoa(*port)
	http.ListenAndServe(addr, r)
}
