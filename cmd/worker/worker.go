package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/client"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("0.0.0.0:8082", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	sigs := make(chan os.Signal, 1)
	go func() {
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		select {
		case s := <-sigs:
			switch s {
			case syscall.SIGINT, syscall.SIGTERM:
				os.Exit(1)
			}
		}
	}()

	c := server.NewFetcherServiceClient(conn)

	w := client.NewFetcherWorker(c)
	w.Start()
}
