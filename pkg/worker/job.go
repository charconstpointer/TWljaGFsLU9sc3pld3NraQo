package worker

import (
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fetcher"
)

//Job is a wrapper around a Measure with cancel chan in case you want to stop execution
type Job struct {
	Done chan struct{}
	M    *fetcher.Measure
}

//NewJob returns new job
func NewJob(m *fetcher.Measure) *Job {
	return &Job{M: m, Done: make(chan struct{}, 1)}
}

//Cancel signals a need to stop job's execution
func (j *Job) Cancel() {
	j.Done <- struct{}{}
}
