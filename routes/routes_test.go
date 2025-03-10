package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/undersleep7x/cryptowallet-v0.1/handlers"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"dummy": "response"}`))
}

func TestPingRoute(t *testing.T) {
	original := handlers.Ping
	defer func() { handlers.Ping = original }()
	handlers.Ping = dummyHandler
	router := mux.NewRouter()
	SetupRoutes(router)

	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ping returned wrong status: %v", rr.Code)
	}
}

func TestPriceRoute(t *testing.T) {
	original := handlers.FetchPrices
	defer func() { handlers.FetchPrices = original }()
	handlers.FetchPrices = dummyHandler
	router := mux.NewRouter()
	SetupRoutes(router)

	req, err := http.NewRequest("GET", "/price", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Price returned wrong status: %v", rr.Code)
	}
}