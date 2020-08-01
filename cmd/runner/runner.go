package main

import (
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fw"
	"github.com/rs/zerolog/log"
)

func main() {
	w := fw.NewWorker()

	jobDesc := fw.NewDesc("https://google.com", 5)
	job := fw.NewJob(jobDesc)
	err := w.AddJob(job)
	if err != nil {
		log.Error().Err(err)
	}
}
