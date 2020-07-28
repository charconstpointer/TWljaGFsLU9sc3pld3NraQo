package measure

//Measures is a repository layer for an measure aggregate
type Measures interface {
	CreateMeasure(m *Measure) error
	GetMeasure(ID int) (*Measure, error)
	GetMeasures() ([]*Measure, error)
	UpdateMeasure(m *Measure) error
	DeleteMeasure(ID int) error
}
