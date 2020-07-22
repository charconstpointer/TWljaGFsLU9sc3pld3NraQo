package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/router"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server"
	"github.com/labstack/gommon/log"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on for incoming http requests")
	flag.Parse()

	srv := server.NewServer()

	r := router.New(srv)

	addr := ":" + strconv.Itoa(*port)
	s := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error(err)
	}
}
