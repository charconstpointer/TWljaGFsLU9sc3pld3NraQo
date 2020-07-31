package worker

import (
	"time"

	"github.com/rs/zerolog/log"
)

//Job represents a task that should be executed on a set interval
type Job struct {
	D chan struct{}
	T *time.Ticker
	u Unit
}

//Result is an outcome of a single job execution
type Result struct {
	ID      int
	URL     string
	Success bool
	Res     string
	Dur     float64
	Date    time.Time
}

//NewJob .
func NewJob(u *Unit) *Job {
	return &Job{
		u: *u,
		D: make(chan struct{}, 1),
		T: time.NewTicker(time.Duration(u.interval) * time.Second),
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

//Unit .
func (j *Job) Unit() Unit {
	return j.u
}
