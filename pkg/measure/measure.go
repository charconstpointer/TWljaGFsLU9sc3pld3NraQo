package measure

//Measure represents set of properties for a worker to fetch and work on
type Measure struct {
	id       int
	url      string
	interval int
	probes   []interface{}
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

//AsDto returns a measure dto
func (m *Measure) AsDto() Dto {
	return Dto{ID: m.id, URL: m.url, Interval: m.interval}
}
