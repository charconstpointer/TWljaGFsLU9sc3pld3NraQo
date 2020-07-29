package fetcher

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
)

type FetchHandlerSuite struct {
	suite.Suite
	f Fetcher
}

func TestFetchHandlerSuite(t *testing.T) {
	suite.Run(t, new(FetchHandlerSuite))
}

func (suite *FetchHandlerSuite) SetupTest() {
	repo := measure.NewMeasuresRepo()
	f := NewFetcher(repo)
	suite.f = *f
}

func (suite *FetchHandlerSuite) TestHandleCreateMeasure() {

	cm := measure.CreateMeasure{
		Interval: 100,
		URL:      "https://foo.bar",
	}
	reqBody, _ := json.Marshal(cm)
	req := httptest.NewRequest("POST", "http://example.com/foo", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	suite.f.HandleCreateMeasure(w, req)

	resp := w.Result()

	// body, _ := ioutil.ReadAll(resp.Body)
	suite.Equal(200, resp.StatusCode)

	cm = measure.CreateMeasure{
		Interval: 100,
		URL:      "rewq",
	}
	reqBody, _ = json.Marshal(cm)
	req = httptest.NewRequest("POST", "http://example.com/foo", bytes.NewBuffer(reqBody))
	w = httptest.NewRecorder()

	suite.f.HandleCreateMeasure(w, req)

	resp = w.Result()

	// body, _ := ioutil.ReadAll(resp.Body)
	suite.Equal(400, resp.StatusCode)

}

func (suite *FetchHandlerSuite) TestHandleDeleteMeasure() {
	req := httptest.NewRequest("DELETE", "http://example.com/api/fetcher/1", nil)
	w := httptest.NewRecorder()

	suite.f.HandleDeleteMeasure(w, req)

	resp := w.Result()

	suite.Equal(400, resp.StatusCode)
}

func (suite *FetchHandlerSuite) TestHandleUpdateMeasure() {
	cm := measure.CreateMeasure{
		Interval: 100,
		URL:      "https://foo.bar",
	}
	reqBody, _ := json.Marshal(cm)

	req := httptest.NewRequest("UPDATE", "http://example.com/foo", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	suite.f.HandleCreateMeasure(w, req)

	resp := w.Result()

	// body, _ := ioutil.ReadAll(resp.Body)
	suite.Equal(200, resp.StatusCode)

	cm = measure.CreateMeasure{
		Interval: 102,
		URL:      "https://foo.bar",
	}
	reqBody, _ = json.Marshal(cm)

	req = httptest.NewRequest("UPDATE", "http://example.com/foo", bytes.NewBuffer(reqBody))
	w = httptest.NewRecorder()

	suite.f.HandleCreateMeasure(w, req)

	resp = w.Result()
	suite.Equal(200, resp.StatusCode)

	cmx := map[string]interface{}{
		"interval": "foo",
	}
	reqBody, _ = json.Marshal(cmx)

	req = httptest.NewRequest("UPDATE", "http://example.com/foo", bytes.NewBuffer(reqBody))
	w = httptest.NewRecorder()

	suite.f.HandleCreateMeasure(w, req)

	resp = w.Result()
	suite.Equal(400, resp.StatusCode)
}

func (suite *FetchHandlerSuite) TestHandleGetAllMeasurements() {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	suite.f.HandleGetAllMeasures(w, req)
	resp := w.Result()
	suite.Equal(404, resp.StatusCode)

	m := measure.NewMeasure("https://foo.bar", 10)
	suite.f.enqueue(m)
	req = httptest.NewRequest("GET", "http://example.com/foo", nil)
	w = httptest.NewRecorder()
	suite.f.HandleGetAllMeasures(w, req)

	resp = w.Result()
	suite.Equal(200, resp.StatusCode)
}
