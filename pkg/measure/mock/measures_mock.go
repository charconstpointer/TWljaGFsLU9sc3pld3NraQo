package mock

import (
	"fmt"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
)

type MeasuresMock struct {
	m []*measure.Measure
}

func NewMeasuresMock() *MeasuresMock {
	m := []*measure.Measure{
		{
			ID:       1,
			URL:      "https://wp.pl",
			Interval: 13,
			Probes:   make([]*measure.Probe, 0),
		},
		{
			ID:       2,
			URL:      "https://google.com",
			Interval: 13,
			Probes:   make([]*measure.Probe, 0),
		},
	}
	return &MeasuresMock{m: m}
}

func (m2 MeasuresMock) Save(m *measure.Measure) (int, error) {
	m2.m = append(m2.m, m)
	return 1, nil
}

func (m2 MeasuresMock) SaveProbe(ID int, p measure.Probe) error {
	return nil
}

func (m2 MeasuresMock) Get(ID int) (*measure.Measure, error) {
	for _, m3 := range m2.m {
		if m3.ID == ID {
			return m3, nil
		}
	}
	return nil, fmt.Errorf("%d could not be found", ID)
}

func (m2 MeasuresMock) GetByUrl(URL string) (*measure.Measure, error) {
	for _, m3 := range m2.m {
		if m3.URL == URL {
			return m3, nil
		}
	}
	return nil, fmt.Errorf("%s could not be found", URL)
}

func (m2 MeasuresMock) GetAll() ([]*measure.Measure, error) {
	return m2.m, nil
}

func (m2 MeasuresMock) Update(ID int, interval int) error {
	return nil
}

func (m2 MeasuresMock) Delete(ID int) error {
	var rid = -1
	for i, m := range m2.m {
		if m.ID == ID {
			rid = i
		}
	}
	if rid == -1 {
		return fmt.Errorf("%d could not be found", ID)
	}
	m2.m = append(m2.m[:rid], m2.m[rid+1:]...)
	return nil
}
