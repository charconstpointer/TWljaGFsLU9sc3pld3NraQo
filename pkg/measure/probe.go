package measure

//Probe represents single measurement of a given measure configuration
type Probe struct {
	response  string
	duration  float32
	createdAt float32
}

//NewProbe creates new probe
func NewProbe(response string, duration float32, createdAt float32) *Probe {
	return &Probe{
		response:  response,
		duration:  duration,
		createdAt: float32(createdAt),
	}
}

//ProbeDto is
type ProbeDto struct {
	Response  string  `json:"response"`
	Duration  float32 `json:"duration"`
	CreatedAt float32 `json:"created_at"`
}

//AsDto converts Probe to ProbeDto
func (p *Probe) AsDto() ProbeDto {
	return ProbeDto{
		Response:  p.response,
		Duration:  p.duration,
		CreatedAt: p.createdAt,
	}
}
