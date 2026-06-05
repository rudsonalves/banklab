package bootstrap

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppToken                  string
	JWTSecret                 string
	JWTAccessTokenDuration    time.Duration
	JWTRefreshTokenDuration   time.Duration
	TransactionPasswordPepper string
	Port                      string
	Database                  DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

const (
	defaultJWTAccessTokenDuration  = 15 * time.Minute
	defaultJWTRefreshTokenDuration = 7 * 24 * time.Hour
)

// Init initializes the application by loading environment variables
// and registering errors.
func Init() {
	loadEnv()
	RegisterErrors()
}

// LoadConfig reads configuration from environment variables and returns a Config struct.
// It performs fail-fast checks to ensure required variables are set, logging a fatal error if any are missing.
func LoadConfig() Config {
	appToken := os.Getenv("APP_TOKEN")
	if appToken == "" {
		log.Fatal("APP_TOKEN environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	jwtAccessTokenDuration := durationEnvOrDefault(
		"JWT_ACCESS_TOKEN_DURATION", defaultJWTAccessTokenDuration)
	jwtRefreshTokenDuration := durationEnvOrDefault(
		"JWT_REFRESH_TOKEN_DURATION", defaultJWTRefreshTokenDuration)

	transactionPasswordPepper := os.Getenv("TRANSACTION_PASSWORD_PEPPER")
	if transactionPasswordPepper == "" {
		log.Fatal("TRANSACTION_PASSWORD_PEPPER environment variable is required")
	}

	if len(transactionPasswordPepper) < 32 {
		log.Fatal("TRANSACTION_PASSWORD_PEPPER must be at least 32 characters")
	}

	if transactionPasswordPepper == appToken {
		log.Fatal("TRANSACTION_PASSWORD_PEPPER must not match APP_TOKEN")
	}

	if transactionPasswordPepper == jwtSecret {
		log.Fatal("TRANSACTION_PASSWORD_PEPPER must not match JWT_SECRET")
	}

	protocol := os.Getenv("SERVER_PROTOCOL")
	if protocol == "" {
		protocol = "http"
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost"
	}

	databaseConfig := DatabaseConfig{
		Host:     requiredEnv("DB_HOST"),
		Port:     requiredEnv("DB_PORT"),
		Name:     requiredEnv("DB_NAME"),
		User:     requiredEnv("DB_USER"),
		Password: requiredEnv("DB_PASSWORD"),
	}

	log.Printf("URL to access the server, defaulting to %s://%s:%s", protocol, host, port)

	return Config{
		AppToken:                  appToken,
		JWTSecret:                 jwtSecret,
		JWTAccessTokenDuration:    jwtAccessTokenDuration,
		JWTRefreshTokenDuration:   jwtRefreshTokenDuration,
		TransactionPasswordPepper: transactionPasswordPepper,
		Port:                      port,
		Database:                  databaseConfig,
	}
}

func requiredEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("%s environment variable is required", name)
	}

	return value
}

func durationEnvOrDefault(name string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}

	return parseDurationEnv(name, value)
}

func parseDurationEnv(name string, value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("%s environment variable must be a valid duration: %v", name, err)
	}
	if duration <= 0 {
		log.Fatalf("%s environment variable must be greater than zero", name)
	}

	return duration
}

// loadEnv attempts to load environment variables from .env files
// in multiple locations.
// It checks the current directory, the executable's directory,
// and the parent of the executable's directory.
// If no .env file is found in these locations, it falls back to
// loading from the default location.
func loadEnv() {
	candidates := []string{".env", filepath.Join("api", ".env")}

	if executablePath, err := os.Executable(); err == nil {
		executableDir := filepath.Dir(executablePath)
		candidates = append(candidates,
			filepath.Join(executableDir, ".env"),
			filepath.Join(filepath.Dir(executableDir), ".env"),
		)
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			_ = godotenv.Load(candidate)
			return
		}
	}

	_ = godotenv.Load()
}
