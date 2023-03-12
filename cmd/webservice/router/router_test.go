package router_test

import (
	"github.com/jattento/passenger-tracker/cmd/webservice/router"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	r := router.New(mockCalculator{}, http.NewServeMux())

	req, err := http.NewRequest(http.MethodPost, "/calculate", strings.NewReader(`[["EZE","MIA"],["MIA","BRA"],["BRA","EZE"],["EZE","NYC"]]`))
	if err != nil {
		t.Fatal(err)
	}

	h, pattern := r.Handler(req)

	if h == nil {
		t.Fatal("no handler found")
	}

	if pattern != "/calculate" {
		t.Fatal("wrong calculate endpoint pattern")
	}
}

func TestPing(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(router.PingHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `pong`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
