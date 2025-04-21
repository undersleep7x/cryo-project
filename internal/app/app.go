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
	postgresInfra "github.com/undersleep7x/cryo-project/internal/infra/postgres"
	platformPostgres "github.com/undersleep7x/cryo-project/internal/platform/postgresstore"
	platformRedis "github.com/undersleep7x/cryo-project/internal/platform/redisstore"
	"github.com/undersleep7x/cryo-project/internal/prices"
	"github.com/undersleep7x/cryo-project/internal/transactions"
)

type App struct {
	Config     *config.AppConfig
	RedisCache platformRedis.RedisClient
	PostgresDB platformPostgres.PostgresClient
	Router     *gin.Engine
}

// load configuration file for implementation
func loadAppConfig() *App {
	cfg := config.LoadConfig()

	log.Println("Initializing logging...")
	setupLogging(cfg)

	log.Println("Initializing Postgres DB...")
	postgresClient := setupPgDatabase(cfg)

	log.Println("Loading Redis cache...")
	redisClient := setupRedisCache(cfg)

	log.Println("Wiring interfaces and router...")
	router := gin.Default()
	priceCache := cacheInfra.NewPriceCache(redisClient)
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
		RedisCache: redisClient,
		Router:     router,
		PostgresDB: postgresClient,
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

func setupPgDatabase(cfg *config.AppConfig) platformPostgres.PostgresClient{
	db, err := postgresInfra.NewPostgresClient(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to open Postgres connection: %v", err)
	}
	pgClient := platformPostgres.NewPgClientWrapper(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pgClient.Ping(ctx); err != nil {
		log.Fatalf("Postgres client ping failed: %v", err)
	}

	log.Println("Postgres connected successfully")
	return pgClient
}

func setupRedisCache(cfg *config.AppConfig) platformRedis.RedisClient {
	redisAddr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	rawRedisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	redisClient := platformRedis.NewRedisClientWrapper(rawRedisClient)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	log.Println("Redis connected successfully")
	return redisClient
}

// startup application and configurations
func InitApp() *App {
	log.Println("Initializing config...")
	app := loadAppConfig()
	log.Println("App initialized")
	return app
}
