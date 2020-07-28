package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/worker"
	"google.golang.org/grpc"
)

func main() {
	fmt.Fprintf(os.Stdout, "worker\n")
	conn, err := grpc.Dial("0.0.0.0:8082", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := server.NewFetcherServiceClient(conn)

	w := worker.NewFetcherWorker(c)
	w.Do()

	time.Sleep(99999999 * time.Second)
}
