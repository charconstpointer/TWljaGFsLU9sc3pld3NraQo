package fw

type Desc struct {
	id       int
	url      string
	interval int
}

func NewDesc(url string, interval int) Desc {
	return Desc{
		url:      url,
		interval: interval,
	}
}
