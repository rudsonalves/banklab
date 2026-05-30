package bootstrap

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	AppToken                  string
	JWTSecret                 string
	TransactionPasswordPepper string
}

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

	return Config{
		AppToken:                  appToken,
		JWTSecret:                 jwtSecret,
		TransactionPasswordPepper: transactionPasswordPepper,
	}
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
