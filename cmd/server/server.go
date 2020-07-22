package main

import (
	"flag"
	"net/http"
	"strconv"

	"github.com/charconstpointer/TWljaGFsLU9sc3pld3NraQo/pkg/router"
)

func main() {
	port := flag.Int("port", 8080, "port to listen on for incoming http requests")
	flag.Parse()

	r := router.New()

	addr := ":" + strconv.Itoa(*port)
	http.ListenAndServe(addr, r)
}
