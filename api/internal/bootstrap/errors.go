package bootstrap

import (
	accountapplication "github.com/seu-usuario/bank-api/internal/account/application"
	adminapplication "github.com/seu-usuario/bank-api/internal/admin/application"
	authapplication "github.com/seu-usuario/bank-api/internal/auth/application"
	customerapplication "github.com/seu-usuario/bank-api/internal/customer/application"
)

// RegisterErrors registers all application errors in a centralized manner.
// This function should be called during application initialization to ensure
// that all errors are registered before they are used.
// It allows for a consistent error handling strategy across the application.
func RegisterErrors() {
	// Register application errors
	// 1. Generic errors
	accountapplication.RegisterErrors()
	adminapplication.RegisterErrors()
	// 2. Domain-specific errors
	customerapplication.RegisterErrors()
	// 3. Auth-specific errors
	authapplication.RegisterErrors()
}
