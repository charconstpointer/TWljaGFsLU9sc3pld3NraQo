package main

import (
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/fw"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	w := fw.NewWorker()
	err := w.Start()
	if err != nil {
		log.Error().Msg(err.Error())
	}

	jobDesc := fw.NewDesc("https://google.com", 5)
	job := fw.NewJob(jobDesc)

	err = w.AddJob(job)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	<-w.Done()
}
