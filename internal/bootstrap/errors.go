package bootstrap

import (
	accountapplication "github.com/seu-usuario/bank-api/internal/account/application"
	authapplication "github.com/seu-usuario/bank-api/internal/auth/application"
	customerapplication "github.com/seu-usuario/bank-api/internal/customer/application"
)

func RegisterErrors() {
	// Register application errors
	// 1. Generic errors
	accountapplication.RegisterErrors()
	// 2. Domain-specific errors
	customerapplication.RegisterErrors()
	// 3. Auth-specific errors
	authapplication.RegisterErrors()
}
