package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockPriceFetcher struct {
	result interface{}
	err error
}

func (m mockPriceFetcher) FetchCryptoPrice(cryptoList []string, currency string) (interface{}, error) {
	return m.result, m.err
}

func TestPing(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Ping)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ping handler returned wrong status: %v", status)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response", err)
	}
	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}
}

func TestFetchPrices_MissingCrypto(t *testing.T) {
	req, err := http.NewRequest("GET", "/price?currency=USD", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FetchPrices)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Did not receive missing crypto 400 status: %v", rr.Code)
	}
}

func TestFetchPrices_MissingCurrency(t *testing.T) {
	req, err := http.NewRequest("GET", "/price?crypto=BTC", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FetchPrices)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Did not receive missing currency 400 status: %v", rr.Code)
	}
}

func TestFetchPrices_Success(t *testing.T) {
	originalFetcher := DefaultPriceFetcher
	defer func() { DefaultPriceFetcher = originalFetcher }()

	DefaultPriceFetcher = mockPriceFetcher {
		result: map[string]float64{"bitcoin": 45000.00, "ethereum": 3000.00},
		err: nil,
	}

	req, err := http.NewRequest("GET", "/price?crypto=bitcoin,ethereum&currency=usd", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FetchPrices)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Did not receive 200 OK: %v", rr.Code)
	}

	var result map[string]float64
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if result["bitcoin"] != 45000.00 {
		t.Errorf("Did not receive expected bitcoin price: %v", result["bitcoin"])
	}
	if result["ethereum"] != 3000.00 {
		t.Errorf("Did not receive expected ethereum price: %v", result["ethereum"])
	}
}

func TestFetchPrices_ServiceError(t *testing.T) {
	originalFetcher := DefaultPriceFetcher
	defer func() { DefaultPriceFetcher = originalFetcher }()

	DefaultPriceFetcher = mockPriceFetcher {
		result: nil,
		err: errors.New("dummy error"),
	}

	req, err := http.NewRequest("GET", "/price?crypto=bitcoin,ethereum&currency=usd", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FetchPrices)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Did not receive expected 500 status: %v", rr.Code)
	}
}