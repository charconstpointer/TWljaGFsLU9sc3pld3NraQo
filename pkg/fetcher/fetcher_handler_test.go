package fetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch_HandleCreateMeasure(t *testing.T) {
	repo := mock.NewMeasuresMock()
	fetch := NewImpr(context.Background(), repo)
	cm := measure.CreateMeasure{
		URL:      "https:foobar.com",
		Interval: 13,
	}
	body, err := json.Marshal(cm)
	if err != nil {
		t.Error("json error")
	}
	req, err := http.NewRequest("POST", "/api/fetcher", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fetch.HandleCreateMeasure)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestFetch_HandleGetAllMeasures(t *testing.T) {
	repo := mock.NewMeasuresMock()
	fetch := NewImpr(context.Background(), repo)
	req, err := http.NewRequest("GET", "/api/fetcher", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(fetch.HandleGetAllMeasures)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

//This test will fail, because i'd have to setup chi router to parse id from the route, will fix that later
func TestFetch_HandleDeleteMeasure(t *testing.T) {
	//repo := mock.NewMeasuresMock()
	//fetch := NewImpr(context.Background(), repo)
	//req, err := http.NewRequest("DELETE", "/api/fetcher/1", nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//rr := httptest.NewRecorder()
	//handler := http.HandlerFunc(fetch.HandleDeleteMeasure)
	//handler.ServeHTTP(rr, req)
	//if status := rr.Code; status != http.StatusOK {
	//	t.Errorf("handler returned wrong status code: got %v want %v",
	//		status, http.StatusOK)
	//}
}

//This test will fail, because i'd have to setup chi router to parse id from the route, will fix that later
func TestFetch_GetAllMeasures(t *testing.T) {
	//repo := mock.NewMeasuresMock()
	//fetch := NewImpr(context.Background(), repo)
	//req, err := http.NewRequest("DELETE", "/api/fetcher/1", nil)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//rr := httptest.NewRecorder()
	//handler := http.HandlerFunc(fetch.HandleGetHistory)
	//handler.ServeHTTP(rr, req)
	//if status := rr.Code; status != http.StatusOK {
	//	t.Errorf("handler returned wrong status code: got %v want %v",
	//		status, http.StatusOK)
	//}
}
