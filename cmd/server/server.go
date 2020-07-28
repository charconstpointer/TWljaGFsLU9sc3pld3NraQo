package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server"
	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/server/router"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on for incoming http requests")
	flag.Parse()

	db, err := sqlx.Connect("mysql", "root:app_root_pass@tcp(localhost:3306)/foo")
	if err != nil {
		log.Error(err)
	}

	srv := server.NewServer(db)
	r := router.New(srv)

	addr := ":" + strconv.Itoa(*port)
	s := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error(err)
	}
}
