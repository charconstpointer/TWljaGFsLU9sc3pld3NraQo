package measure

import (
	"fmt"
	"sync"
)

//MeasuresRepo is an implementation of a Measures interface
type MeasuresRepo struct {
	m  []*Measure
	mu sync.RWMutex
}

//NewMeasuresRepo is
func NewMeasuresRepo() *MeasuresRepo {
	return &MeasuresRepo{}
}

//Save persists new Measure
func (msr *MeasuresRepo) Save(m *Measure) error {
	// msr.mu.Lock()
	// defer msr.mu.Unlock()

	msr.m = append(msr.m, m)

	return nil
}

//Delete deletes existing Measure
func (msr *MeasuresRepo) Delete(ID int) error {
	// msr.mu.Lock()
	// defer msr.mu.Unlock()

	i, _ := msr.find(ID)
	if i == -1 {
		return fmt.Errorf("measure with id %d could not be found", ID)
	}

	msr.m = append(msr.m[:i], msr.m[i+1:]...)
	return nil
}

//Get is
func (msr *MeasuresRepo) Get(ID int) (*Measure, error) {
	msr.mu.RLock()
	defer msr.mu.RUnlock()
	i, m := msr.find(ID)
	if i == -1 {
		return nil, fmt.Errorf("measure with id %d could not be found", ID)
	}
	return m, nil
}

//GetAll returns all active measures
func (msr *MeasuresRepo) GetAll() ([]*Measure, error) {
	msr.mu.RLock()
	defer msr.mu.RUnlock()
	return msr.m, nil
}

//Update is
func (msr *MeasuresRepo) Update(ID int, interval int) error {
	msr.mu.Lock()
	defer msr.mu.Unlock()
	for _, msr := range msr.m {
		if msr.id == ID {
			msr.interval = interval
			return nil
		}
	}
	return fmt.Errorf("measure with id %d could not be found", ID)
}

//GetByUrl is
func (msr *MeasuresRepo) GetByUrl(URL string) (*Measure, error) {
	msr.mu.RLock()
	defer msr.mu.RUnlock()
	for _, m := range msr.m {
		if m.url == URL {
			return m, nil
		}
	}
	return nil, fmt.Errorf("measure for url %s could not be found", URL)
}

func (msr *MeasuresRepo) find(ID int) (int, *Measure) {
	msr.mu.RLock()
	defer msr.mu.RUnlock()
	for i, m := range msr.m {
		if m.id == ID {
			return i, m
		}
	}
	return -1, nil
}
