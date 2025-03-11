package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/undersleep7x/cryptowallet-v0.1/handlers"
)

//setup dummy response for calls to handler
func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"dummy": "response"}`))
}

func TestPingRoute(t *testing.T) {
	//replace original handler with dummy handler, and set defer to restore original during test cleanup
	original := handlers.Ping
	defer func() { handlers.Ping = original }()
	handlers.Ping = dummyHandler
	router := mux.NewRouter()
	SetupRoutes(router)

	//send request to ping endpoint and ensure response is received
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	//ping should return ok
	if rr.Code != http.StatusOK {
		t.Errorf("Ping returned wrong status: %v", rr.Code)
	}
}

func TestPriceRoute(t *testing.T) {
	//replace original handler with dummy handler, and set defer to restore original during test cleanup
	original := handlers.FetchPrices
	defer func() { handlers.FetchPrices = original }()
	handlers.FetchPrices = dummyHandler
	router := mux.NewRouter()
	SetupRoutes(router)

	//send request to price endpoint and ensure response is received
	req, err := http.NewRequest("GET", "/price", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	//dummy handler should return ok as long as router works properly
	if rr.Code != http.StatusOK {
		t.Errorf("Price returned wrong status: %v", rr.Code)
	}
}