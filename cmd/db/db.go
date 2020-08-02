package main

import (
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/measure"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type meas struct {
}

func main() {
	db, _ := sqlx.Connect("mysql", "root:password@tcp(127.0.0.1:3306)/foobar")
	repo := measure.MySQLRepo{
		DB: db,
	}

	_, _ = repo.Get(1)
	repo.Save(measure.NewMeasure("https://latozradiem.pl", 1000))
	repo.Save(measure.NewMeasure("https://latozradiem.pl", 5000))
}
