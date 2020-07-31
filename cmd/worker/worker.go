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

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var (
	grpcPort    = flag.Int("grpc", 8082, "server grpc port")
	timeout     = flag.Int("timeout", 5000, "http client timeout in ms")
	grpcTimeout = flag.Int("dial", 5000, "grpc dial timeout in ms")
)

func main() {
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	g := errgroup.Group{}
	g.Go(func() error {
		log.Info().Msg("starting gRPC connection with a server")

		conn, err := grpc.Dial(
			fmt.Sprintf(":%d", *grpcPort),
			grpc.WithInsecure(), grpc.WithBlock(),
			grpc.WithTimeout(time.Duration(*grpcTimeout)*time.Millisecond),
		)
		defer conn.Close()

		if err != nil {
			log.Fatal().Msg("can't connect to fetcher server")
			return err
		}

		log.Info().Msg("connected to fetcher server")

		c := fetcher.NewFetcherServiceClient(conn)
		repo := worker.NewProbesRepo(c)
		w := worker.NewWorker(*timeout, repo)
		log.Info().Msg("starting worker")
		err = w.Start(ctx)

		return err
	})

	select {
	case <-interrupt:
		cancel()
		break
	}

	err := g.Wait()
	if err != nil {
		log.Error().Err(err)
		os.Exit(2)
	}

}
