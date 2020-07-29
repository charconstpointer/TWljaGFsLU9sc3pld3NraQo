package measure

//Measures is a repository layer for an measure aggregate
type Measures interface {
	Save(m *Measure) error
	Get(ID int) (*Measure, error)
	Exists(URL string) bool
	GetAll() ([]*Measure, error)
	Update(m *Measure) error
	Delete(ID int) error
}
