package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/undersleep7x/cryptowallet-v0.1/app"
	"github.com/undersleep7x/cryptowallet-v0.1/services/api"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockAPI struct {
	mock.Mock
}

func (m *MockAPI) FetchPrices(cryptoList []string, currency string) (*resty.Response, error) {
	args := m.Called(cryptoList, currency)
	return args.Get(0).(*resty.Response), args.Error(1)
}
func TestFetchCryptoPrice(t *testing.T) {
	cryptoSymbols := []string{"bitcoin"}
	currency := "usd"

	t.Run("Cache Hit", func(t *testing.T) {
		mockRedis := new(MockRedisClient)
		app.RedisClient = mockRedis

		cachedData := `{"bitcoin": 45000.00}`
		mockRedis.On("Get", mock.Anything, "prices:bitcoin:usd").Return(cachedData, nil)

		prices, err := FetchCryptoPrice(cryptoSymbols, currency)
		assert.NoError(t, err)
		assert.Equal(t, 45000.00, prices["bitcoin"])
	})
	
	t.Run("Cache Miss - API Failure", func(t *testing.T) {
		mockRedis := new(MockRedisClient)
		mockAPI := new(MockAPI)
		app.RedisClient = mockRedis
	
		originalFetchPrices := api.FetchPrices
		defer func() { api.FetchPrices = originalFetchPrices }()
		api.FetchPrices = func(cryptoList []string, currency string) (*resty.Response, error) {
			return mockAPI.FetchPrices(cryptoList, currency)
		}
	
		mockRedis.On("Get", mock.Anything, "prices:bitcoin:usd").Return("", errors.New("redis: nil"))
	
		dummyErrorResponse := &resty.Response{}
		dummyErrorResponse.SetBody([]byte (`{}`))
	
		mockAPI.On("FetchPrices", mock.Anything, mock.Anything).Return(dummyErrorResponse, errors.New("API error"))
	
		mockRedis.On("Set", mock.Anything, "prices:bitcoin:usd", mock.Anything, 30*time.Second).Return(nil).Maybe()
	
		prices, err := FetchCryptoPrice(cryptoSymbols, currency)
		assert.NoError(t, err)
		assert.Equal(t, -1.00, prices["bitcoin"])
	})

	t.Run("Redis Error", func(t *testing.T) {
		mockRedis := new(MockRedisClient)
		app.RedisClient = mockRedis
	
		mockRedis.On("Get", mock.Anything, "prices:bitcoin:usd").Return("", errors.New("redis connection error"))
	
		mockAPI := new(MockAPI)
		originalFetchPrices := api.FetchPrices
		defer func() { api.FetchPrices = originalFetchPrices }()
		api.FetchPrices = func(cryptoList []string, currency string) (*resty.Response, error) {
			return mockAPI.FetchPrices(cryptoList, currency)
		}
	
		dummyResponse := &resty.Response{}
		dummyResponse.SetBody([]byte(`{"bitcoin":{"usd":47000.00}}`))
	
		mockAPI.On("FetchPrices", mock.Anything, mock.Anything).Return(dummyResponse, nil)
		mockRedis.On("Set", mock.Anything, "prices:bitcoin:usd", mock.Anything, 30*time.Second).Return(nil)
	
		prices, err := FetchCryptoPrice(cryptoSymbols, currency)
		assert.NoError(t, err)
		assert.Equal(t, 47000.00, prices["bitcoin"])
	})

	t.Run("Cache Miss - API Success", func(t *testing.T) {
		mockRedis := new(MockRedisClient)
		mockAPI := new(MockAPI)
		app.RedisClient = mockRedis

		originalFetchPrices := api.FetchPrices
		defer func() { api.FetchPrices = originalFetchPrices }()
		api.FetchPrices = func(cryptoList []string, currency string) (*resty.Response, error) {
			return mockAPI.FetchPrices(cryptoList, currency)
		}

		mockRedis.On("Get", mock.Anything, "prices:bitcoin:usd").Return("", errors.New("redis: nil"))

		dummyResponse := &resty.Response{}
		dummyResponse.SetBody([]byte (`{"bitcoin":{"usd":46000.00}}`))

		mockAPI.On("FetchPrices", mock.Anything, mock.Anything).Return(dummyResponse, nil)
		mockRedis.On("Set", mock.Anything, "prices:bitcoin:usd", mock.Anything, 30*time.Second).Return(nil)

		prices, err := FetchCryptoPrice(cryptoSymbols, currency)
		assert.NoError(t, err)
		assert.Equal(t, 46000.00, prices["bitcoin"])
	})
}
