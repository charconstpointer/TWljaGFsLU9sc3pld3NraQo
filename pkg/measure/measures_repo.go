package measure

import "fmt"

//MeasuresRepo is an implementation of a Measures interface
type MeasuresRepo struct {
	m []*Measure
}

//NewMeasuresRepo is
func NewMeasuresRepo() *MeasuresRepo {
	return &MeasuresRepo{}
}

//CreateMeasure persists new Measure
func (msr *MeasuresRepo) CreateMeasure(c CreateMeasure) error {
	measure := c.AsEntity()
	msr.m = append(msr.m, measure)
	fmt.Printf("Creating new measure %v", c)
	return nil
}

//DeleteMeasure deletes existing Measure
func (msr *MeasuresRepo) DeleteMeasure(ID int) error {
	fmt.Printf("Deleting measure %d", ID)
	return nil
}

//GetMeasure is
func (msr *MeasuresRepo) GetMeasure(ID int) (*Measure, error) {
	return nil, nil
}

//GetMeasures returns all active measures
func (msr *MeasuresRepo) GetMeasures() ([]*Measure, error) {
	return msr.m, nil
}

//UpdateMeasure is
func (msr *MeasuresRepo) UpdateMeasure(c CreateMeasure) error {
	return nil
}
