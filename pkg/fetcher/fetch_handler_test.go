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
}

func TestFetchHandlerSuite(t *testing.T) {
	suite.Run(t, new(FetchHandlerSuite))
}

func (suite *FetchHandlerSuite) TestHandleCreateMeasure() {
	repo := measure.NewMeasuresRepo()
	f := NewFetcher(repo)
	cm := measure.CreateMeasure{
		Interval: 100,
		URL:      "https://foo.bar",
	}
	reqBody, _ := json.Marshal(cm)
	req := httptest.NewRequest("POST", "http://example.com/foo", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	f.HandleCreateMeasure(w, req)

	resp := w.Result()

	// body, _ := ioutil.ReadAll(resp.Body)
	suite.Equal(200, resp.StatusCode)
}

func (suite *FetchHandlerSuite) TestHandleDeleteMeasure() {
	repo := measure.NewMeasuresRepo()
	f := NewFetcher(repo)

	req := httptest.NewRequest("DELETE", "http://example.com/api/fetcher/1", nil)
	w := httptest.NewRecorder()

	f.HandleDeleteMeasure(w, req)

	resp := w.Result()

	// body, _ := ioutil.ReadAll(resp.Body)
	suite.Equal(400, resp.StatusCode)
}
