package fetcher

import (
	"context"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure/mock"
	"testing"
)

func TestFetch_CreateOrUpdate(t *testing.T) {
	repo := mock.NewMeasuresMock()
	fetch := NewImpr(context.Background(), repo)
	cm := measure.CreateMeasure{
		URL:      "https://foobar.com",
		Interval: 12,
	}
	id, err := fetch.CreateOrUpdate(cm)
	if err != nil {
		t.Error(err.Error())
	}

	if id != 3 {
		t.Errorf("wrong id returned, expected 3, got %d", id)
	}

	cm = measure.CreateMeasure{
		URL:      "https://foobar.com",
		Interval: 15,
	}
	id, err = fetch.CreateOrUpdate(cm)
	if err != nil {
		t.Error(err.Error())
	}
	if id != 3 {
		t.Errorf("id should remain unchanged, got %d", id)
	}
}

func TestFetch_DeleteMeasure(t *testing.T) {
	repo := mock.NewMeasuresMock()
	fetch := NewImpr(context.Background(), repo)
	id := 1

	err := fetch.DeleteMeasure(id)
	if err != nil {
		t.Error(err.Error())
	}
	m, err := fetch.GetAllMeasures()
	if err != nil {
		t.Error(err.Error())
	}
	for _, ms := range m {
		if ms.ID == id {
			t.Errorf("measure with id %d was not deleted", id)
		}
	}

	err = fetch.DeleteMeasure(id)
	if err == nil {
		t.Errorf("err is nil, expected an error because measure with id %d is no longer present", id)
	}
}

func TestFetch_GetAllMeasures2(t *testing.T) {
	repo := mock.NewMeasuresMock()
	fetch := NewImpr(context.Background(), repo)
	m, err := fetch.GetAllMeasures()
	if err != nil {
		t.Error(err.Error())
	}
	if len(m) != 2 {
		t.Errorf("incorrect amount of measures was returned expected 2, got %d", len(m))
	}
}

func TestFetch_GetHistory(t *testing.T) {
	repo := mock.NewMeasuresMock()
	fetch := NewImpr(context.Background(), repo)
	p, err := fetch.GetHistory(1)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(p) != 3 {
		t.Errorf("incorrect amount of probes was returned expected 3, got %d", len(p))
	}
}
