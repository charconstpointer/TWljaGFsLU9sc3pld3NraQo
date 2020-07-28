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
func (msr *MeasuresRepo) CreateMeasure(c CreateMeasure) error {
	msr.mu.Lock()
	defer msr.mu.Unlock()

	measure := c.AsEntity()
	msr.m = append(msr.m, measure)

	return nil
}

//DeleteMeasure deletes existing Measure
func (msr *MeasuresRepo) DeleteMeasure(ID int) error {
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
func (msr *MeasuresRepo) GetMeasure(ID int) (*Measure, error) {
	i, m := msr.find(ID)
	if i == -1 {
		return nil, fmt.Errorf("measure with id %d could not be found", ID)
	}
	return m, nil
}

//GetMeasures returns all active measures
func (msr *MeasuresRepo) GetMeasures() ([]*Measure, error) {
	return msr.m, nil
}

//UpdateMeasure is
func (msr *MeasuresRepo) UpdateMeasure(c CreateMeasure) error {
	return nil
}

func (msr *MeasuresRepo) find(ID int) (int, *Measure) {
	for i, m := range msr.m {
		if m.id == ID {
			return i, m
		}
	}
	return -1, nil
}
