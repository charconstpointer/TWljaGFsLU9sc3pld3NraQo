package worker

import "time"

type Result struct {
	ID      int
	URL     string
	Success bool
	Res     string
	Dur     float64
	Date    time.Time
}
