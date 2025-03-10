package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

type ConfigStruct struct { //utilize config files to start app services
	App struct {
		Name string `yaml:"name"`
		Environment string `yaml:"environment"`
		LogLevel string `yaml:"log_level"`
		Port string `yaml:"port"`
	} `yaml:"app"` //create app struct

	Logging struct {
		FilePath   string `yaml:"file_path"`
		MaxSize    int    `yaml:"max_size"`
		MaxBackups int    `yaml:"max_backups"`
		MaxAge     int    `yaml:"max_age"`
		Permissions int `yaml:"permissions"`
	} `yaml:"logging"` //create logging struct

	Redis struct {
		Address string `yaml:"address"`
		Password string `yaml:"password"`
		DB int `yaml:"db"`
		Timeout int `yaml:"timeout"`
		Ttl int `yaml:"ttl"`
	} `yaml:"redis"` //create redis cache struct

	CoinGecko struct {
		BaseURL string `yaml:"base_url"`
		Timeout int `yaml:"timeout"`
		RetryAttempts int `yaml:"retry_attempts"`
	} `yaml:"coingecko"` //create coingecko api struct
}

var Config ConfigStruct
var RedisClient *redis.Client
var Router *mux.Router

func loadConfig() {
	env := os.Getenv("APP_ENV")
	if env == ""{
		log.Println("No application environment found; loading with local env config.")
		env = "local"
	}

	configFile := fmt.Sprintf("config/%s_config.yml", env)

	log.Printf("Loading config: %s", configFile)
	data, err := os.ReadFile("utils/config/local_config.yml")
	if err != nil {
		log.Fatalf("Failed to read config file %v", err)
	}
	err = yaml.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}
	log.Println("Config initialized")
}

func setupRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: Config.Redis.Address,
		Password: Config.Redis.Password,
		DB: Config.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	} else {
		log.Println("Redis connected successfully")
	}
}

func setupLogging() {
	logFile, err := os.OpenFile(Config.Logging.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(Config.Logging.Permissions))
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

func InitApp() {
	log.Println("Initializing config...")
	loadConfig()
	log.Println("Initializing logging...")
	setupLogging()
	log.Println("Loading Redis cache...")
	setupRedis()
	log.Println("Setting router...")
	Router = mux.NewRouter()
	log.Println("App initialized")
}