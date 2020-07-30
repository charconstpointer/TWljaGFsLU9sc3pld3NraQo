package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/worker"
	"github.com/labstack/gommon/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
	grpcPort = flag.Int("grpc", 8082, "server grpc port")
	timeout  = flag.Int("timeout", 5000, "http client timeout in ms")
)

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	g := errgroup.Group{}
	g.Go(func() error {
		log.Print("starting gRPC connection with a server\n")
		conn, err := grpc.Dial(fmt.Sprintf(":%d", *grpcPort), grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
		defer conn.Close()

		if err != nil {
			return err
		}
		log.Print("connected")
		c := fetcher.NewFetcherServiceClient(conn)
		repo := worker.NewProbesRepo(c)

		w := worker.NewWorker(repo)
		log.Print("starting worker")
		err = w.Start(ctx)
		return err
	})

	select {
	case <-interrupt:
		cancel()
		// os.Exit(2)
		break
	}

	err := g.Wait()
	if err != nil {
		log.Error(err)
		os.Exit(2)
	}

}
