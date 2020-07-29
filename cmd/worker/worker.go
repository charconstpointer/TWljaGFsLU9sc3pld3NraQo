package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/worker"
	"github.com/labstack/gommon/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
	grpcPort = flag.Int("grpc", 8082, "server grpc port")
	timeout  = flag.Int("timeout", 5000, "http client timeout in ms")

	grpcServer *grpc.Server
)

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)
	g := errgroup.Group{}
	g.Go(func() error {
		log.Print("Starting new grcp conn \n")
		conn, err := grpc.Dial(fmt.Sprintf(":%d", *grpcPort), grpc.WithInsecure(), grpc.WithBlock())
		defer conn.Close()

		if err != nil {
			return err
		}

		c := fetcher.NewFetcherServiceClient(conn)
		w := worker.NewFetcherWorker(c, *timeout)

		return w.Start()
	})

	select {
	case <-interrupt:
		os.Exit(2)
		break
	}

	if grpcServer != nil {
		grpcServer.GracefulStop()
	}

	err := g.Wait()
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}

}
