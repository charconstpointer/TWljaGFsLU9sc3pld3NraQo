package main

import (
	"flag"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher/router"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	"google.golang.org/grpc"

	"github.com/labstack/gommon/log"
)

var (
	httpPort := flag.Int("http", 8080, "port to listen on for incoming http requests")
	grpcPort := flag.Int("grpc", 8082, "port to listen on for incoming grpc requests")
)

func main() {
	flag.Parse()

	
	

	repo := measure.NewMeasuresRepo()
	srv := fetcher.NewFetcher(repo)
	r := router.New(srv)

	s := &http.Server{
		Addr:    httpAddr,
		Handler: r,
	}

	go func() {
		grpcAddr := ":" + strconv.Itoa(*grpcPort)
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Errorf("failed to listen: %v", err)
		}

		gs := grpc.NewServer()

		fetcher.RegisterFetcherServiceServer(gs, srv)
		log.Infof("Starting gRPC server %v", time.Now())
		gs.Serve(lis)
	}()

	go func() {
		log.Infof("Starting http server %v", time.Now())
	httpAddr := ":" + strconv.Itoa(*httpPort)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error(err)
	}
	}()

}
