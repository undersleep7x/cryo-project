package prices

import (
	"context"
	"errors"
	"testing"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/undersleep7x/cryo-project/internal/infra/cache"
)

// mock redis and api call for testing
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Mock.Called(ctx, key)
	return args.String(0), args.Error(1)
}
func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Mock.Called(ctx, key, value, expiration)
	return args.Error(0)
}
func (m *MockRedisClient) Ping(ctx context.Context) error {
	args := m.Mock.Called(ctx)
	return args.Error(0)
}

type MockAPI struct {
	mock.Mock
}

func (m *MockAPI) FetchPrices(cryptoList []string, currency string, baseURL string, timeoutVal int) (*resty.Response, error) {
	args := m.Mock.Called(cryptoList, currency, baseURL, timeoutVal)
	return args.Get(0).(*resty.Response), args.Error(1)
}

func TestFetchCryptoPrice(t *testing.T) {
	// set starter variables for test cases
	cryptoSymbols := []string{"bitcoin"}
	currency := "usd"
	testConfig := Config{
		BaseURL:       "https://dummy-coingecko.com",
		Timeout:       5,
		RetryAttempts: 1,
	}

	// test response if value is found in cache
	t.Run("Cache Hit", func(t *testing.T) {
		// set mock redis cache and test data
		mockRedis := new(MockRedisClient)
		mockPriceCache := cache.NewPriceCache(mockRedis)
		service := NewFetchCryptoPriceService(mockPriceCache, testConfig)

		cachedData := `{"bitcoin": 45000.00}`
		mockRedis.Mock.On("Get", mock.Anything, "prices:bitcoin:usd").Return(cachedData, nil)

		// make method call and record response, should have no error and match test data
		prices, err := service.FetchCryptoPrice(cryptoSymbols, currency)
		assert.NoError(t, err)
		assert.Equal(t, 45000.00, prices["bitcoin"])
	})

	// check and fail api after cache failure
	t.Run("Cache Miss - API Failure", func(t *testing.T) {
		// set mock redis and api call
		mockAPI := new(MockAPI)
		mockRedis := new(MockRedisClient)
		mockPriceCache := cache.NewPriceCache(mockRedis)
		service := NewFetchCryptoPriceService(mockPriceCache, testConfig)

		// switch the method for the mock method and revert after ending test
		originalFetchPrices := FetchPrices
		defer func() { FetchPrices = originalFetchPrices }()
		FetchPrices = func(cryptoList []string, currency string, baseURL string, timeoutVal int) (*resty.Response, error) {
			return mockAPI.FetchPrices(cryptoList, currency, baseURL, timeoutVal)
		}

		// set mock redis behavior to trigger api check
		mockRedis.Mock.On("Get", mock.Anything, "prices:bitcoin:usd").Return("", errors.New("redis connection error"))

		// set dummy response from resty call in method
		dummyErrorResponse := &resty.Response{}
		dummyErrorResponse.SetBody([]byte(`{}`))

		// set mock api behavior and post call redis behavior
		mockAPI.Mock.On("FetchPrices", cryptoSymbols, currency, testConfig.BaseURL, testConfig.Timeout).Return(dummyErrorResponse, nil)
		mockRedis.Mock.On("Set", mock.Anything, "prices:bitcoin:usd", mock.Anything, 30*time.Second).Return(nil)

		// make method call, should return fallback price for crypto
		prices, err := service.FetchCryptoPrice(cryptoSymbols, currency)
		assert.NoError(t, err)
		assert.Equal(t, -1.00, prices["bitcoin"])
	})

	t.Run("Redis Error", func(t *testing.T) {
		// set mock redis and api
		mockAPI := new(MockAPI)
		mockRedis := new(MockRedisClient)
		mockPriceCache := cache.NewPriceCache(mockRedis)
		service := NewFetchCryptoPriceService(mockPriceCache, testConfig)

		// set mock redis response
		mockRedis.Mock.On("Get", mock.Anything, "prices:bitcoin:usd").Return("", errors.New("redis connection error"))

		// switch the method for the mock method and revert after ending test
		originalFetchPrices := FetchPrices
		defer func() { FetchPrices = originalFetchPrices }()
		FetchPrices = func(cryptoList []string, currency string, baseURL string, timeoutVal int) (*resty.Response, error) {
			return mockAPI.FetchPrices(cryptoList, currency, baseURL, timeoutVal)
		}

		//set dummy resty response and mock api response/ post api call redis action
		dummyResponse := &resty.Response{}
		dummyResponse.SetBody([]byte(`{"bitcoin":{"usd":47000.00}}`))
		mockAPI.Mock.On("FetchPrices", cryptoSymbols, currency, testConfig.BaseURL, testConfig.Timeout).Return(dummyResponse, nil)
		mockRedis.Mock.On("Set", mock.Anything, "prices:bitcoin:usd", mock.Anything, 30*time.Second).Return(nil)

		// make method call, should return expected price for crypto
		prices, err := service.FetchCryptoPrice(cryptoSymbols, currency)
		assert.NoError(t, err)
		assert.Equal(t, 47000.00, prices["bitcoin"])
	})

	t.Run("Cache Miss - API Success", func(t *testing.T) {
		// set mock redis and api
		mockAPI := new(MockAPI)
		mockRedis := new(MockRedisClient)
		mockPriceCache := cache.NewPriceCache(mockRedis)
		service := NewFetchCryptoPriceService(mockPriceCache, testConfig)

		// switch the method for the mock method and revert after ending test
		originalFetchPrices := FetchPrices
		defer func() { FetchPrices = originalFetchPrices }()
		FetchPrices = func(cryptoList []string, currency string, baseURL string, timeoutVal int) (*resty.Response, error) {
			return mockAPI.FetchPrices(cryptoList, currency, baseURL, timeoutVal)
		}

		// set mock responses from redis and api call
		mockRedis.Mock.On("Get", mock.Anything, "prices:bitcoin:usd").Return("", errors.New("redis: nil"))
		dummyResponse := &resty.Response{}
		dummyResponse.SetBody([]byte(`{"bitcoin":{"usd":46000.00}}`))
		mockAPI.Mock.On("FetchPrices", cryptoSymbols, currency, testConfig.BaseURL, testConfig.Timeout).Return(dummyResponse, nil)
		mockRedis.Mock.On("Set", mock.Anything, "prices:bitcoin:usd", mock.Anything, 30*time.Second).Return(nil)

		// make method call, should fail to find in redis and return from api call
		prices, err := service.FetchCryptoPrice(cryptoSymbols, currency)
		assert.NoError(t, err)
		assert.Equal(t, 46000.00, prices["bitcoin"])
	})
}
