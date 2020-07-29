package measure

import (
	"fmt"
	"sync"
)

//MeasuresRepo is an implementation of a Measures interface
type MeasuresRepo struct {
	m  []*Measure
	mu sync.Mutex
}

//NewMeasuresRepo is
func NewMeasuresRepo() *MeasuresRepo {
	return &MeasuresRepo{}
}

//CreateMeasure persists new Measure
func (msr *MeasuresRepo) Save(m *Measure) error {
	msr.mu.Lock()
	defer msr.mu.Unlock()

	msr.m = append(msr.m, m)

	return nil
}

//DeleteMeasure deletes existing Measure
func (msr *MeasuresRepo) Delete(ID int) error {
	msr.mu.Lock()
	defer msr.mu.Unlock()

	i, _ := msr.find(ID)
	if i == -1 {
		return fmt.Errorf("measure with id %d could not be found", ID)
	}

	msr.m = append(msr.m[:i], msr.m[i+1:]...)
	return nil
}

//GetMeasure is
func (msr *MeasuresRepo) Get(ID int) (*Measure, error) {
	i, m := msr.find(ID)
	if i == -1 {
		return nil, fmt.Errorf("measure with id %d could not be found", ID)
	}
	return m, nil
}

//GetMeasures returns all active measures
func (msr *MeasuresRepo) GetAll() ([]*Measure, error) {
	return msr.m, nil
}

//UpdateMeasure is
func (msr *MeasuresRepo) Update(m *Measure) error {
	for _, msr := range msr.m {
		if msr.url == m.url {
			msr.interval = m.interval
			return nil
		}
	}
	return fmt.Errorf("measure %s could not be found", m.url)
}
func (msr *MeasuresRepo) Exists(URL string) bool {
	for _, m := range msr.m {
		if m.url == URL {
			return true
		}
	}
	return false
}

func (msr *MeasuresRepo) find(ID int) (int, *Measure) {
	for i, m := range msr.m {
		if m.id == ID {
			return i, m
		}
	}
	return -1, nil
}
