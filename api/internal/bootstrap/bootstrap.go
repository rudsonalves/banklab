package bootstrap

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func Init() {
	loadEnv()
	RegisterErrors()
}

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
