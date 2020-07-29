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
	httpPort := flag.Int("http", 8080, "port to listen on for incoming http requests")
	grpcPort := flag.Int("grpc", 8082, "port to listen on for incoming grpc requests")

	flag.Parse()

	httpAddr := ":" + strconv.Itoa(*httpPort)
	grpcAddr := ":" + strconv.Itoa(*grpcPort)

	repo := measure.NewMeasuresRepo()
	srv := server.NewServer(repo)
	r := router.New(srv)

	s := &http.Server{
		Addr:    httpAddr,
		Handler: r,
	}

	log.Infof("Starting http server %v", time.Now())
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error(err)
	}

	go func() {

		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Errorf("failed to listen: %v", err)
		}

		gs := grpc.NewServer()

		server.RegisterFetcherServiceServer(gs, srv)
		log.Infof("Starting gRPC server %v", time.Now())
		gs.Serve(lis)
	}()
}
