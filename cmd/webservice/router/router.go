package router

import (
	"log"
	"net/http"
)

// Router works as manager for all the handlers
type Router struct {
	*http.ServeMux
	calculator Calculator
}

// Calculator for allowing dependency injection
type Calculator interface {
	FirstAndLast(flights ...[]string) (firstAirport, lastAirport string, err error)
}

// New returns a new Router
func New(calc Calculator, serv *http.ServeMux) *Router {
	initializeRouter(serv, calc)

	return &Router{
		calculator: calc,
		ServeMux:   serv,
	}
}

func initializeRouter(router *http.ServeMux, calc Calculator) {
	router.HandleFunc("/calculate", CalculatorHandlerGetFirstAndLastAirport(calc))

	// Add a ping router just for testing...
	router.HandleFunc("/ping", PingHandler)
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong"))
	log.Println("pong")
}

// generic error message model
type errorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
