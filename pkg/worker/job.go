package worker

import (
	"time"

	"github.com/rs/zerolog/log"
)

type job interface {
	Cancel()
}

//Job represents a task that should be executed on a set interval
type Job struct {
	D chan (struct{})
	T *time.Ticker
	p Probe
}

//Result is an outcome of a single job execution
type Result struct {
	Probe   int
	URL     string
	Success bool
	Res     string
	Dur     float64
	Date    time.Time
}

//NewJob .
func NewJob(p *Probe) *Job {
	return &Job{
		p: *p,
		D: make(chan struct{}, 1),
		T: time.NewTicker(time.Duration(p.interval) * time.Second),
	}
}

//Cancel sends singal to terminate job exection
func (j *Job) Cancel() {
	select {
	case j.D <- struct{}{}:
		log.Info().Msgf("cancelling job %v", &j)
	default:
		log.Info().Msgf("can not cancell job %v", &j)
	}
}

//Probe .
func (j *Job) Probe() Probe {
	return j.p
}
