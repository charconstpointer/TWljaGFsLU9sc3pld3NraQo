package measure

import (
	"time"
)

//Probe represents single measurement of a given measure configuration
type Probe struct {
	response  string
	duration  float64
	createdAt time.Time
}

//NewProbe creates new probe
func NewProbe(response string, duration float64) *Probe {
	return &Probe{
		response:  response,
		duration:  duration,
		createdAt: time.Now(),
	}
}

//ProbeDto is
type ProbeDto struct {
	Response  string    `json:"response"`
	Duration  float64   `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
}

//AsDto converts Probe to ProbeDto
func (p *Probe) AsDto() ProbeDto {
	return ProbeDto{
		Response:  p.response,
		Duration:  p.duration,
		CreatedAt: p.createdAt,
	}
}
