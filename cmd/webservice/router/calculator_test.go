package router_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jattento/passenger-tracker/cmd/webservice/router"
	"github.com/jattento/passenger-tracker/internal/calculator"
)

func TestCalculatorHandlerGetFirstAndLastAirport_OK(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/calculate", strings.NewReader(`[["EZE","MIA"],["MIA","BRA"],["BRA","EZE"],["EZE","NYC"]]`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(router.CalculatorHandlerGetFirstAndLastAirport(mockCalculator{FirstAndLastFn: func(flights ...[]string) (string, string, error) {
		if slicesOfSlicesEqual(flights, [][]string{{"EZE", "MIA"}, {"MIA", "BRA"}, {"BRA", "EZE"}, {"EZE", "NYC"}}) {
			return "EZE", "NYC", nil
		}
		return "", "", errors.New("unexpected input at mock calculator")
	}}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `["EZE","NYC"]
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCalculatorHandlerGetFirstAndLastAirport_ERR(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/calculate", strings.NewReader(`[["EZE","MIA"],["MIA","BRA"],["BRA","EZE"],["NIC","NYC"]]`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(router.CalculatorHandlerGetFirstAndLastAirport(mockCalculator{FirstAndLastFn: func(flights ...[]string) (string, string, error) {
		if slicesOfSlicesEqual(flights, [][]string{{"EZE", "MIA"}, {"MIA", "BRA"}, {"BRA", "EZE"}, {"NIC", "NYC"}}) {
			return "", "", calculator.ErrFlightsNotConnected
		}
		return "", "", errors.New("unexpected input at mock calculator")
	}}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"code":400,"message":"not all flights are connected"}
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

type mockCalculator struct {
	FirstAndLastFn func(...[]string) (string, string, error)
}

func (mock mockCalculator) FirstAndLast(flights ...[]string) (firstAirport, lastAirport string, err error) {
	return mock.FirstAndLastFn(flights...)
}

func slicesOfSlicesEqual(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, s := range a {
		found := false
		for _, t := range b {
			if slicesEqual(s, t) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aMap := make(map[string]int)
	bMap := make(map[string]int)
	for _, v := range a {
		aMap[v]++
	}
	for _, v := range b {
		bMap[v]++
	}
	for k, v := range aMap {
		if bMap[k] != v {
			return false
		}
	}
	return true
}
