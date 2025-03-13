package routes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupTestRouter() *gin.Engine {
	router := gin.Default()
	SetupRoutes(router)
	return router
}

func PerformRequest(router *gin.Engine, method, path string, body []byte) *httptest.ResponseRecorder{
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestPingRoute(t *testing.T) {
	router := SetupTestRouter()
	w := PerformRequest(router, "GET", "/", nil)
	assert.NotEqual(t, http.StatusNotFound, w.Code)
}

func TestPriceRoute(t *testing.T) {
	router := SetupTestRouter()
	w := PerformRequest(router, "GET", "/price", nil)
	assert.NotEqual(t, http.StatusNotFound, w.Code)
}
