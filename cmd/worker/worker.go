package main

import (
	"context"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	figure.NewColorFigure("worker", "slant", "blue", true).Print()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := grpc.Dial(
		fmt.Sprintf(":%d", 8084),
		grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithTimeout(time.Duration(5000)*time.Millisecond),
	)
	defer conn.Close()

	log.Info().Msg("✅ connected to fetcher server")

	c := fetcher.NewFetcherServiceClient(conn)
	bp := worker.NewFetcherBackplane(c)

	w := worker.NewWorker(bp)
	go w.Start(ctx)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	select {
	case <-interrupt:
		err := w.Stop()
		if err != nil {
			log.Error().Msg(err.Error())
		}
		log.Info().Msg("interrupt")
		os.Exit(0)
	}

}
