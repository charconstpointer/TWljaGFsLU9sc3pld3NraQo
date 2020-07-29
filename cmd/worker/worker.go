package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/worker"
	"google.golang.org/grpc"
)

var (
	grpcPort = flag.Int("grpc", 8082, "server grpc port")

	grpcServer *grpc.Server
)

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	go func() {
		conn, err := grpc.Dial(fmt.Sprintf(":%d", *grpcPort), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := fetcher.NewFetcherServiceClient(conn)

		w := worker.NewFetcherWorker(c)
		w.Start()
	}()

	select {
	case <-interrupt:
		break
	}

	if grpcServer != nil {
		grpcServer.GracefulStop()
	}
}
