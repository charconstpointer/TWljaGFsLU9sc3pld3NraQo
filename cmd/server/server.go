package main

import (
	"flag"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server/router"
	"google.golang.org/grpc"

	"github.com/labstack/gommon/log"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on for incoming http requests")
	flag.Parse()

	repo := measure.NewMeasuresRepo()
	srv := server.NewServer(repo)

	r := router.New(srv)

	addr := ":" + strconv.Itoa(*port)
	s := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go func() {

		lis, err := net.Listen("tcp", "0.0.0.0:8082")
		if err != nil {
			log.Errorf("failed to listen: %v", err)
		}

		gs := grpc.NewServer()

		server.RegisterFetcherServiceServer(gs, srv)
		log.Infof("Starting gRPC server %v", time.Now())
		gs.Serve(lis)
	}()

	log.Infof("Starting http server %v", time.Now())
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error(err)
	}

}
