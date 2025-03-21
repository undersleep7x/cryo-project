package app

import (
	"context"
	"os"
	"testing"
	"time"
)

// dummy config file for testing
func createDummyConfigFile(t *testing.T) {
	configDir := "config"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	dummyConfigContent := `app:
  name: "TestWallet"
  environment: "local"
  log_level: "debug"
  port: "8080"
logging:
  file_path: "test_log.log"
  max_size: 10
  max_backups: 3
  max_age: 30
  permissions: 0644
redis:
  address: "localhost:6379"
  password: ""
  db: 0
  timeout: 5
  ttl: 300
coingecko:
  base_url: "https://api.coingecko.com"
  timeout: 5
  retry_attempts: 3
`
	configFilePath := "config/local_config.yml"
	if err := os.WriteFile(configFilePath, []byte(dummyConfigContent), 0644); err != nil {
		t.Fatalf("Failed to write dummy config file: %v", err)
	}
}

// clear the creating of the dummy log file when spinning up the config
func removeDummyConfigFile(t *testing.T) {
	configFilePath := "config/"
	if err := os.RemoveAll(configFilePath); err != nil {
		t.Fatalf("Failed to remove dummy config file: %v", err)
	}
}

func TestLoadConfig(t *testing.T) {
	// setup pre-testing configurations with logs and config file
	os.Setenv("APP_ENV", "local")
	defer os.Unsetenv("APP_ENV")
	createDummyConfigFile(t)
	defer removeDummyConfigFile(t)

	InitApp()

	// check that config exists and that its on a port
	if Config.App.Name == "" {
		t.Fatal("Expected config to load; received empty App.Name")
	}
	if Config.App.Port == "" {
		t.Fatal("Expected app port to be set; received empty value")
	}
}

func TestSetupRedis(t *testing.T) {
	// config file setup
	createDummyConfigFile(t)
	defer removeDummyConfigFile(t)

	// set initial redis configurations for testing
	Config.Redis.Address = "localhost:6379"
	Config.Redis.Password = ""
	Config.Redis.DB = 1

	InitApp() // start app and redis server context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// redis ping should return a response, not an error, if online
	err := RedisClient.Ping(ctx)
	if err != nil {
		t.Fatalf("Redis connection failed: %v", err)
	}
}

func TestSetupLogging(t *testing.T) {
	// config file setup
	createDummyConfigFile(t)
	defer removeDummyConfigFile(t)

	// set logging configs for test
	tempLogFile := "test_log.log"
	Config.Logging.FilePath = tempLogFile
	Config.Logging.Permissions = 0644

	InitApp() //start app

	// templogfile should exists with no errors
	_, err := os.Stat(tempLogFile)
	if os.IsNotExist(err) {
		t.Fatalf("Expected log file to be created: %s", tempLogFile)
	}
	os.Remove(tempLogFile)
}
