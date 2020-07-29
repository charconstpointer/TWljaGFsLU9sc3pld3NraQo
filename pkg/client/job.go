package client

import "github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server"

//Job is a wrapper around a Measure with cancel chan in case you want to stop execution
type Job struct {
	Done chan struct{}
	M    *server.Measure
}

//NewJob returns new job
func NewJob(m *server.Measure) *Job {
	return &Job{M: m, Done: make(chan struct{}, 1)}
}

//Cancel signals a need to stop job's execution
func (j *Job) Cancel() {
	j.Done <- struct{}{}
}
