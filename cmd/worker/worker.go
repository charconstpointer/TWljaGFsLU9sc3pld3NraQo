package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server"
	"google.golang.org/grpc"
)

func main() {
	fmt.Fprintf(os.Stdout, "worker\n")
	conn, err := grpc.Dial("0.0.0.0:8082", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	c := server.NewFetcherServiceClient(conn)
	msr, err := c.GetMeasures(ctx, &server.GetMeasuresRequest{})
	if err != nil {
		fmt.Println(err)
	}

	for _, m := range msr.Measures {
		go func(m *server.Measure, fc server.FetcherServiceClient) {
			fmt.Printf("Loaded measure : %v\n", m)
			t := time.NewTicker(time.Duration(m.Interval) * time.Second)
			c := http.Client{
				Timeout: time.Duration(m.Interval) * time.Second,
			}
			for {
				select {
				case _ = <-t.C:
					fmt.Printf("Fetching : %s\n", m.URL)
					var start time.Time
					var duration int64
					req, err := http.NewRequest("GET", m.URL, nil)
					tracer := &httptrace.ClientTrace{
						ConnectStart: func(network, addr string) {
							start = time.Now()
						},
						ConnectDone: func(network, addr string, err error) {
							duration = time.Since(start).Milliseconds()
						},
					}
					req = req.WithContext(httptrace.WithClientTrace(req.Context(), tracer))
					res, err := c.Do(req)

					if err != nil {
						fmt.Println(err)
					}
					b, err := ioutil.ReadAll(res.Body)
					if err != nil {
						fmt.Println(err)
					}

					ctx, cancel := context.WithTimeout(context.Background(), time.Second)

					defer cancel()
					fc.AddProbe(ctx, &server.AddProbeRequest{
						MeasureID: m.ID,
						CreatedAt: time.Now().UnixNano(),
						Duration:  float32(float64(duration) / float64(time.Millisecond)),
						Response:  string(b),
					})

				}
			}
		}(m, c)
	}
	time.Sleep(99999999 * time.Second)
}
