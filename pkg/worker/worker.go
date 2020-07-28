package worker

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server"
)

type worker interface {
	Listen()
}

type FetcherWorker struct {
	c server.FetcherServiceClient
}

func NewFetcherWorker(c server.FetcherServiceClient) *FetcherWorker {
	return &FetcherWorker{c}
}

func (w *FetcherWorker) Listen() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msr, err := w.c.GetMeasures(ctx, &server.GetMeasuresRequest{})
	if err != nil {
		return err
	}

	for _, m := range msr.Measures {
		go func(m *server.Measure, fc server.FetcherServiceClient, ctx context.Context) {
			log.Printf("Loaded measure : %v\n", m)
			t := time.NewTicker(time.Duration(m.Interval) * time.Second)
			for {
				select {
				case _ = <-t.C:
					fetch(m, w.c)
				case _ = <-ctx.Done():
					return
				}
			}
		}(m, w.c, context.Background())
	}

	return nil
}

func fetch(m *server.Measure, fc server.FetcherServiceClient) {
	log.Printf("Fetching : %s\n", m.URL)

	var start time.Time
	var duration int64

	req, err := http.NewRequest("GET", m.URL, nil)
	c := http.Client{
		Timeout: time.Duration(m.Interval) * time.Second,
	}

	start = time.Now()
	res, err := c.Do(req)
	duration = time.Since(start).Milliseconds()

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
