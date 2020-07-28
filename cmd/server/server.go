package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server/router"
  
	"github.com/labstack/gommon/log"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on for incoming http requests")
	flag.Parse()
  
	repo := measure.NewMeasuresRepo()
	srv := server.NewServer(repo)


	srv := server.NewServer(db)
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
