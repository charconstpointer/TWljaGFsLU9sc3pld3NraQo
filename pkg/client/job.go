package client

import "github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server"

type Job struct {
	Done chan struct{}
	M    *server.Measure
}

func NewJob(m *server.Measure) *Job {
	return &Job{M: m, Done: make(chan struct{}, 1)}
}

func (j *Job) Cancel() {
	j.Done <- struct{}{}
}
