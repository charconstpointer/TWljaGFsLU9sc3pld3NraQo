package worker

import (
	"log"
	"time"
)

type job interface {
	Cancel()
}

//Job .
type Job struct {
	D chan (struct{})
	T *time.Ticker
	p Probe
}

//Result .
type Result struct {
	Probe   int
	URL     string
	Success bool
	Res     string
	Dur     int
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

//Cancel .
func (j *Job) Cancel() {
	select {
	case j.D <- struct{}{}:
		log.Printf("cancelling job %v", &j)
	default:
		log.Printf("can not cancell job %v", &j)
	}
}

//Probe .
func (j *Job) Probe() Probe {
	return j.p
}
