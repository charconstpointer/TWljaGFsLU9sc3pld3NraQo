package worker

import (
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
	"testing"
)

func TestWorker_enqueue(t *testing.T) {
	repo := ProbesRepoMock{}
	w := NewWorker(1, repo)
	ev := fetcher.ListenForChangesResponse{
		MeasureID: 1,
		Change:    fetcher.Change_CREATED,
		Measure: &fetcher.Measure{
			ID:       1,
			URL:      "",
			Interval: 1,
		},
	}
	go w.enqueue(&ev)
	select {
	case j, ok := <-w.jobs:
		if j.Probe().id != 1 {
			t.Error("job not added")
		}
		if !ok {
			return
		}

	}

}

func TestWorker_replace(t *testing.T) {
	repo := ProbesRepoMock{}
	w := NewWorker(1, repo)

	ev := fetcher.ListenForChangesResponse{
		MeasureID: 1,
		Change:    fetcher.Change_CREATED,
		Measure: &fetcher.Measure{
			ID:       1,
			URL:      "",
			Interval: 50,
		},
	}
	w.enqueue(&ev)
	ev = fetcher.ListenForChangesResponse{
		MeasureID: 1,
		Change:    fetcher.Change_EDITED,
		Measure: &fetcher.Measure{
			ID:       1,
			URL:      "",
			Interval: 5,
		},
	}
	w.replace(&ev)
	select {
	case job := <-w.jobs:
		if job.Probe().interval != 5 {
			t.Error("..")
		}
	}

}
