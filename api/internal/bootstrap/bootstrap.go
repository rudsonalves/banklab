package bootstrap

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment                  string
	PublicBaseURL                string
	AppToken                     string
	JWTSecret                    string
	JWTAccessTokenDuration       time.Duration
	JWTRefreshTokenDuration      time.Duration
	TransactionPasswordPepper    string
	ExposeDebugVerificationToken bool
	Host                         string
	Port                         string
	Database                     DatabaseConfig
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

const (
	EnvironmentDev        = "dev"
	EnvironmentStaging    = "staging"
	EnvironmentProduction = "production"
)

// Init initializes the application from the explicitly selected environment
// file and registers application errors.
func Init() {
	loadEnv()
	RegisterErrors()
}

// LoadConfig reads configuration from environment variables and returns a Config struct.
// It performs fail-fast checks to ensure required variables are set, logging a fatal error if any are missing.
func LoadConfig() Config {
	environment, err := parseEnvironment(requiredEnv("APP_ENV"))
	if err != nil {
		log.Fatal(err)
	}

	publicBaseURL := requiredEnv("PUBLIC_BASE_URL")

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

	exposeDebugVerificationToken, err := parseRequiredBoolEnv("EXPOSE_DEBUG_VERIFICATION_TOKEN")
	if err != nil {
		log.Fatal(err)
	}
	if err := validateDebugTokenExposure(environment, exposeDebugVerificationToken); err != nil {
		log.Fatal(err)
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

	log.Printf("environment=%s server_address=%s:%s", environment, host, port)

	return Config{
		Environment:                  environment,
		PublicBaseURL:                publicBaseURL,
		AppToken:                     appToken,
		JWTSecret:                    jwtSecret,
		JWTAccessTokenDuration:       jwtAccessTokenDuration,
		JWTRefreshTokenDuration:      jwtRefreshTokenDuration,
		TransactionPasswordPepper:    transactionPasswordPepper,
		ExposeDebugVerificationToken: exposeDebugVerificationToken,
		Host:                         host,
		Port:                         port,
		Database:                     databaseConfig,
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

func parseEnvironment(value string) (string, error) {
	environment := strings.ToLower(strings.TrimSpace(value))
	switch environment {
	case EnvironmentDev, EnvironmentStaging, EnvironmentProduction:
		return environment, nil
	default:
		return "", fmt.Errorf(
			"APP_ENV must be one of %q, %q, or %q",
			EnvironmentDev,
			EnvironmentStaging,
			EnvironmentProduction,
		)
	}
}

func parseRequiredBoolEnv(name string) (bool, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return false, fmt.Errorf("%s environment variable is required", name)
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s environment variable must be true or false", name)
	}

	return parsed, nil
}

func validateDebugTokenExposure(environment string, expose bool) error {
	if environment == EnvironmentProduction && expose {
		return fmt.Errorf("EXPOSE_DEBUG_VERIFICATION_TOKEN must be false in production")
	}

	return nil
}

// loadEnv loads exactly the file selected through ENV_FILE. It deliberately
// avoids fallback discovery so one environment cannot silently use another.
func loadEnv() {
	envFile := strings.TrimSpace(os.Getenv("ENV_FILE"))
	if envFile == "" {
		log.Fatal("ENV_FILE environment variable is required")
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("failed to load ENV_FILE %q: %v", envFile, err)
	}
}
