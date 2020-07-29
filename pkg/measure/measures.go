package measure

//Measures is a repository layer for an measure aggregate
type Measures interface {
	Save(m *Measure) error
	Get(ID int) (*Measure, error)
	GetByUrl(URL string) (*Measure, error)
	GetAll() ([]*Measure, error)
	Update(ID int, interval int) error
	Delete(ID int) error
}
