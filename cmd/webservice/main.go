package main

import (
	"log"
	"net/http"

	"github.com/jattento/passenger-tracker/cmd/webservice/router"
	"github.com/jattento/passenger-tracker/internal/calculator"
)

func main() {
	server := router.New(
		&calculator.Calculator{},
		http.NewServeMux(),
	)

	if err := http.ListenAndServe(":8080", server.ServeMux); err != nil {
		log.Fatal("server error:", err.Error())
	}
}
