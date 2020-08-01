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
	jobs := []fw.Desc{
		fw.NewDesc("https://latozradiem.pl/api/stars", 7),
		fw.NewDesc("https://google.com", 5),
		fw.NewDesc("https://apipodcasts.polskieradio.pl/api/podcasts", 9),
	}

	for _, jd := range jobs {
		job := fw.NewJob(jd)
		_ = w.AddJob(job)
	}

	if err != nil {
		log.Error().Msg(err.Error())
	}

	<-w.Done()
}
