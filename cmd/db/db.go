package main

import (
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type meas struct {
}

func main() {
	db, _ := sqlx.Connect("mysql", "root:password@tcp(127.0.0.1:3306)/foobar")
	repo := measure.MySQLRepo{
		DB: db,
	}
	//
	//_, _ = repo.Get(1)
	//repo.Save(measure.NewMeasure("https://latozradiem.pl", 1000))
	//repo.Save(measure.NewMeasure("https://latozradiem.pl", 5000))
	//m, err := repo.Get(10)
	//if err != nil {
	//	log.Error().Msg(err.Error())
	//	return
	//}
	//log.Info().Msgf("%d", len(m.Probes()))
	ms, err := repo.GetAll()
	if err != nil {
		log.Error().Msg(err.Error())
	}
	log.Info().Int("c", len(ms))
}
