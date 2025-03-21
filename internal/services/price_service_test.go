package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/undersleep7x/cryo-project/internal/app"
)

// mock redis and api call for testing
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
	// set starter variables for test cases
	cryptoSymbols := []string{"bitcoin"}
	currency := "usd"
	service := NewFetchCryptoPriceService()

	// test response if value is found in cache
	t.Run("Cache Hit", func(t *testing.T) {
		// set mock redis cache and test data
		mockRedis := new(MockRedisClient)
		app.RedisClient = mockRedis
		cachedData := `{"bitcoin": 45000.00}`
		mockRedis.On("Get", mock.Anything, "prices:bitcoin:usd").Return(cachedData, nil)

		// make method call and record response, should have no error and match test data
		prices, err := service.FetchCryptoPrice(cryptoSymbols, currency)
		assert.NoError(t, err)
		assert.Equal(t, 45000.00, prices["bitcoin"])
	})

	// check and fail api after cache failure
	t.Run("Cache Miss - API Failure", func(t *testing.T) {
		// set mock redis and api call
		mockRedis := new(MockRedisClient)
		mockAPI := new(MockAPI)
		app.RedisClient = mockRedis

		// switch the method for the mock method and revert after ending test
		originalFetchPrices := FetchPrices
		defer func() { FetchPrices = originalFetchPrices }()
		FetchPrices = func(cryptoList []string, currency string) (*resty.Response, error) {
			return mockAPI.FetchPrices(cryptoList, currency)
		}

		// set mock redis behavior to trigger api check
		mockRedis.On("Get", mock.Anything, "prices:bitcoin:usd").Return("", errors.New("redis: nil"))

		// set dummy response from resty call in method
		dummyErrorResponse := &resty.Response{}
		dummyErrorResponse.SetBody([]byte(`{}`))

		// set mock api behavior and post call redis behavior
		mockAPI.On("FetchPrices", mock.Anything, mock.Anything).Return(dummyErrorResponse, errors.New("API error"))
		mockRedis.On("Set", mock.Anything, "prices:bitcoin:usd", mock.Anything, 30*time.Second).Return(nil).Maybe()

		// make method call, should return fallback price for crypto
		prices, err := service.FetchCryptoPrice(cryptoSymbols, currency)
		assert.NoError(t, err)
		assert.Equal(t, -1.00, prices["bitcoin"])
	})

	t.Run("Redis Error", func(t *testing.T) {
		// set mock redis and api
		mockRedis := new(MockRedisClient)
		app.RedisClient = mockRedis
		mockAPI := new(MockAPI)

		// set mock redis response
		mockRedis.On("Get", mock.Anything, "prices:bitcoin:usd").Return("", errors.New("redis connection error"))

		// switch the method for the mock method and revert after ending test
		originalFetchPrices := FetchPrices
		defer func() { FetchPrices = originalFetchPrices }()
		FetchPrices = func(cryptoList []string, currency string) (*resty.Response, error) {
			return mockAPI.FetchPrices(cryptoList, currency)
		}

		//set dummy resty response and mock api response/ post api call redis action
		dummyResponse := &resty.Response{}
		dummyResponse.SetBody([]byte(`{"bitcoin":{"usd":47000.00}}`))
		mockAPI.On("FetchPrices", mock.Anything, mock.Anything).Return(dummyResponse, nil)
		mockRedis.On("Set", mock.Anything, "prices:bitcoin:usd", mock.Anything, 30*time.Second).Return(nil)

		// make method call, should return expected price for crypto
		prices, err := service.FetchCryptoPrice(cryptoSymbols, currency)
		assert.NoError(t, err)
		assert.Equal(t, 47000.00, prices["bitcoin"])
	})

	t.Run("Cache Miss - API Success", func(t *testing.T) {
		// set mock redis and api
		mockRedis := new(MockRedisClient)
		mockAPI := new(MockAPI)
		app.RedisClient = mockRedis

		// switch the method for the mock method and revert after ending test
		originalFetchPrices := FetchPrices
		defer func() { FetchPrices = originalFetchPrices }()
		FetchPrices = func(cryptoList []string, currency string) (*resty.Response, error) {
			return mockAPI.FetchPrices(cryptoList, currency)
		}

		// set mock responses from redis and api call
		mockRedis.On("Get", mock.Anything, "prices:bitcoin:usd").Return("", errors.New("redis: nil"))
		dummyResponse := &resty.Response{}
		dummyResponse.SetBody([]byte(`{"bitcoin":{"usd":46000.00}}`))
		mockAPI.On("FetchPrices", mock.Anything, mock.Anything).Return(dummyResponse, nil)
		mockRedis.On("Set", mock.Anything, "prices:bitcoin:usd", mock.Anything, 30*time.Second).Return(nil)

		// make method call, should fail to find in redis and return from api call
		prices, err := service.FetchCryptoPrice(cryptoSymbols, currency)
		assert.NoError(t, err)
		assert.Equal(t, 46000.00, prices["bitcoin"])
	})
}
