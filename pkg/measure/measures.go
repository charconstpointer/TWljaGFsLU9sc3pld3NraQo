package measure

//Measures is a repository layer for an measure aggregate
type Measures interface {
	CreateMeasure(c CreateMeasure) error
	GetMeasure(ID int) (*Measure, error)
	GetMeasures() ([]*Measure, error)
	UpdateMeasure(c CreateMeasure) error
	DeleteMeasure(ID int) error
}
