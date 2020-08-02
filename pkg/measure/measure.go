package measure

//Measure represents set of properties for a worker to fetch and work on
type Measure struct {
	id       int
	url      string
	interval int
	probes   []*Probe
}

//Measures is a repository layer for an measure aggregate
type Measures interface {
	Save(m *Measure) (int, error)
	SaveProbe(ID int, p Probe) error
	Get(ID int) (*Measure, error)
	GetByUrl(URL string) (*Measure, error)
	GetAll() ([]*Measure, error)
	Update(ID int, interval int) error
	Delete(ID int) error
}

//NewMeasure is
func NewMeasure(url string, interval int) *Measure {
	return &Measure{
		url:      url,
		interval: interval,
	}
}

//AddProbe adds new probe
func (m *Measure) AddProbe(p *Probe) {
	m.probes = append(m.probes, p)
}

//CreateMeasure represents model that user has to provide in order to create new measure
type CreateMeasure struct {
	URL      string `json:"url"`
	Interval int    `json:"interval"`
}

//Dto is
type Dto struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	Interval int    `json:"interval"`
}

//AsEntity converts CreateMeasure request to a domain entity
func (c CreateMeasure) AsEntity() *Measure {
	return &Measure{
		url:      c.URL,
		interval: c.Interval,
	}
}

//Probes returns probes for a given measure
func (m *Measure) Probes() []*Probe {
	return m.probes
}

//AsDto returns a measure dto
func (m *Measure) AsDto() Dto {
	return Dto{ID: m.id, URL: m.url, Interval: m.interval}
}
