package prices

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockPriceHandler struct {
	FetchCryptoPriceService
}

func (m *mockPriceHandler) FetchCryptoPrice(cryptoList []string, currency string) (map[string]float64, error) {
	if currency == "se" {
		return nil, errors.New("No currency provided")
	}
	return map[string]float64{"bitcoin": 45000.000, "ethereum": 3200.75}, nil
}

func TestFetchPrices(t *testing.T) {
	router := gin.Default()
	mockService := &mockPriceHandler{}
	PriceHandler := NewPriceHandler(mockService)
	router.GET("/price", PriceHandler.FetchPrices)

	t.Run("Missing crypto", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/price?currency=usd", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response map[string]any;
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Missing 'crypto' query parameter", response["error"])
	})

	t.Run("Missing Currency", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/price?crypto=bitcoin,ethereum", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response map[string]any;
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "Missing 'currency' query parameter", response["error"])
	})

	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/price?crypto=bitcoin,ethereum&currency=usd", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	
		assert.Equal(t, http.StatusOK, w.Code)
	
		var response map[string]any;
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}
	
		assert.NotNil(t, response["prices"])
		assert.Equal(t, 45000.00, response["prices"].(map[string]any)["bitcoin"])
		assert.Equal(t, 3200.75, response["prices"].(map[string]any)["ethereum"])
	})

	t.Run("ServiceError", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/price?crypto=bitcoin,ethereum&currency=se", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	
		var response map[string]any;
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}
	
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "Failed to fetch prices", response["error"])
	})
}