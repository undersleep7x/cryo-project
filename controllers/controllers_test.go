package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

//mocking pricefetcher and implementing fetchCryptoPrice due to interface flexibility
type mockPriceFetcher struct {
	result any
	err error
}
func (m mockPriceFetcher) FetchCryptoPrice(cryptoList []string, currency string) (interface{}, error) {
	return m.result, m.err
}

func TestPing(t *testing.T) {
	// send http request to handler for ping to simulate call
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Ping)
	handler.ServeHTTP(rr, req)

	// verify that status is 200
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ping handler returned wrong status: %v", status)
	}

	// verify that response is present and a status of ok
	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response", err)
	}
	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", response["status"])
	}
}

func TestFetchPrices_MissingCrypto(t *testing.T) {
	// test request with missing crypto type and record the response
	req, err := http.NewRequest("GET", "/price?currency=USD", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FetchPrices)
	handler.ServeHTTP(rr, req)

	// response should just be a bad request
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Did not receive missing crypto 400 status: %v", rr.Code)
	}
}

func TestFetchPrices_MissingCurrency(t *testing.T) {
	// test request with missing currency and record the response
	req, err := http.NewRequest("GET", "/price?crypto=BTC", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FetchPrices)
	handler.ServeHTTP(rr, req)

	// response should just be a bad request
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Did not receive missing currency 400 status: %v", rr.Code)
	}
}

func TestFetchPrices_Success(t *testing.T) {
	// swap the original pricefetcher to the mock fetcher and set it to return moc data
	// undo the swap at the end
	originalFetcher := DefaultPriceFetcher
	defer func() { DefaultPriceFetcher = originalFetcher }()
	DefaultPriceFetcher = mockPriceFetcher {
		result: map[string]float64{"bitcoin": 45000.00, "ethereum": 3000.00},
		err: nil,
	}

	// send request with expected crypto and currency, and record response
	req, err := http.NewRequest("GET", "/price?crypto=bitcoin,ethereum&currency=usd", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FetchPrices)
	handler.ServeHTTP(rr, req)

	// status should be 200
	if rr.Code != http.StatusOK {
		t.Errorf("Did not receive 200 OK: %v", rr.Code)
	}

	// result should have expected cryptos and prices
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
	// swap the original pricefetcher to the mock fetcher and set it to return moc data
	// undo the swap at the end
	originalFetcher := DefaultPriceFetcher
	defer func() { DefaultPriceFetcher = originalFetcher }()
	DefaultPriceFetcher = mockPriceFetcher {
		result: nil,
		err: errors.New("dummy error"),
	}

	// send request with expected crypto and currency, and record response
	req, err := http.NewRequest("GET", "/price?crypto=bitcoin,ethereum&currency=usd", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(FetchPrices)
	handler.ServeHTTP(rr, req)

	// response should just be 500 server error
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Did not receive expected 500 status: %v", rr.Code)
	}
}