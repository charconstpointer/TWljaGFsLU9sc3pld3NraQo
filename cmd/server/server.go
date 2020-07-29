package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher/router"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/labstack/gommon/log"
)

var (
	httpPort = flag.Int("http", 8080, "port to listen on for incoming http requests")
	grpcPort = flag.Int("grpc", 8082, "port to listen on for incoming grpc requests")

	grpcServer *grpc.Server
	httpServer *http.Server
)

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	repo := measure.NewMeasuresRepo()
	srv := fetcher.NewFetcher(repo)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		grpcAddr := ":" + strconv.Itoa(*grpcPort)
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			log.Errorf("failed to listen: %v", err)
		}

		grpcServer = grpc.NewServer()

		fetcher.RegisterFetcherServiceServer(grpcServer, srv)
		log.Infof("Starting gRPC server %v", time.Now())
		return grpcServer.Serve(lis)
	})

	g.Go(func() error {
		log.Infof("Starting http server %v", time.Now())
		httpAddr := ":" + strconv.Itoa(*httpPort)
		r := router.New(srv)

		httpServer = &http.Server{
			Addr:    httpAddr,
			Handler: r,
		}
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	select {
	case <-interrupt:
		break
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if httpServer != nil {
		_ = httpServer.Shutdown(shutdownCtx)
	}
	if grpcServer != nil {
		grpcServer.GracefulStop()
	}

	err := g.Wait()
	if err != nil {
		os.Exit(2)
	}

}
