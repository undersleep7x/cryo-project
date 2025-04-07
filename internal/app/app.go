package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
	"github.com/undersleep7x/cryo-project/api/routes"
	"github.com/undersleep7x/cryo-project/internal/config"
	cacheInfra "github.com/undersleep7x/cryo-project/internal/infra/cache"
	platformRedis "github.com/undersleep7x/cryo-project/internal/platform/redisstore"
	"github.com/undersleep7x/cryo-project/internal/prices"
	"github.com/undersleep7x/cryo-project/internal/transactions"
)

type App struct {
	Config     *config.AppConfig
	RedisCache platformRedis.RedisClient
	Router     *gin.Engine
}

// load configuration file for implementation
func loadAppConfig() *App {
	cfg := config.LoadConfig()

	log.Println("Initializing logging...")
	setupLogging(cfg)

	log.Println("Loading Redis cache...")
	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	rawRedisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	RedisClient := platformRedis.NewRedisClientWrapper(rawRedisClient)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := RedisClient.Ping(ctx); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	log.Println("Redis connected successfully")

	log.Println("Wiring interfaces and router...")
	router := gin.Default()
	priceCache := cacheInfra.NewPriceCache(RedisClient)
	priceConfig := prices.Config{
		BaseURL:       "https://api.coingecko.com/api/v3",
		Timeout:       5,
		RetryAttempts: 3,
	}
	priceService := prices.NewFetchCryptoPriceService(priceCache, priceConfig)
	priceHandler := prices.NewPriceHandler(priceService)

	txnRepository := transactions.NewTxnRepository()
	txnService := transactions.NewTransactionsService(txnRepository)
	txnHandler := transactions.NewTransactionsHandler(txnService)
	routes.SetupRoutes(router, priceHandler, txnHandler)

	log.Println("Config initialized")

	return &App{
		Config:     cfg,
		RedisCache: RedisClient,
		Router:     router,
	}
}

// setup logging with logging file
func setupLogging(cfg *config.AppConfig) {
	logFile, err := os.OpenFile(cfg.LoggingPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(0666))
	//opens specified file for logging form config, setting it to be created/appended and read/write only with proper permissions
	if err != nil {
		//fallback & local logging option
		fmt.Printf("Failed to open log file: %v. Logging to stdout.\n", err)
		log.SetOutput(os.Stdout)
	} else {
		//logs to log file & stdout if log file is found
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) //log formatting
	log.Println("Logger initialized")
}

// startup application and configurations
func InitApp() *App {
	log.Println("Initializing config...")
	app := loadAppConfig()
	log.Println("App initialized")
	return app
}
