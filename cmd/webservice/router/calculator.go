package router

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/jattento/passenger-tracker/internal/calculator"
)

func CalculatorHandlerGetFirstAndLastAirport(calc Calculator) func(http.ResponseWriter, *http.Request) {
	// Actual handler implementation...
	return func(w http.ResponseWriter, r *http.Request) {
		// Panic check
		defer func() {
			if err := recover(); err != nil {
				handleErr(w, http.StatusInternalServerError, "Internal error")
			}
		}()
		w.Header().Set("Content-Type", "application/json")

		// POST is the only supported method
		if r.Method != http.MethodPost {
			handleErr(w, http.StatusMethodNotAllowed, "POST is the only supported method on this endpoint")
			return
		}

		flights := make([][]string, 0)

		if err := json.NewDecoder(r.Body).Decode(&flights); err != nil {
			handleErr(w, http.StatusBadRequest, "wrong input format")
			return
		}

		first, last, err := calc.FirstAndLast(flights...)
		if err != nil {
			// Check which type of error we have to give a better user experience...
			if errors.Is(err, calculator.ErrMissingData) {
				handleErr(w, http.StatusBadRequest, "incomplete information")
				return
			}
			if errors.Is(err, calculator.ErrFlightsNotConnected) {
				handleErr(w, http.StatusBadRequest, "not all flights are connected")
				return
			}
			if errors.Is(err, calculator.ErrNondeterministic) {
				handleErr(w, http.StatusBadRequest, "start and end airport are the same, it is not possible to determine with the current information which is the starting")
				return
			}

			handleErr(w, http.StatusInternalServerError, "Internal error")
			return
		}

		logIfError(json.NewEncoder(w).Encode([]string{first, last}))

		// Log input and output for debugging
		log.Println("calculatorHandlerGetFirstAndLastAirport", "input:", flights, "output:", first, last)
	}
}

func handleErr(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	logIfError(json.NewEncoder(w).Encode(errorMessage{code, message}))
}

func logIfError(err error) {
	if err != nil {
		log.Println("ERROR:", err.Error(), "stack:", string(debug.Stack()))
	}
}
