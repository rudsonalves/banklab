package main

import (
	"log"
	"net/http"
	"os"
	"time"

	accountApplication "github.com/seu-usuario/bank-api/internal/account/application"
	accountDelivery "github.com/seu-usuario/bank-api/internal/account/delivery"
	accountInfrastructure "github.com/seu-usuario/bank-api/internal/account/infrastructure"
	authApplication "github.com/seu-usuario/bank-api/internal/auth/application"
	authDelivery "github.com/seu-usuario/bank-api/internal/auth/delivery"
	authInfrastructure "github.com/seu-usuario/bank-api/internal/auth/infrastructure"
	"github.com/seu-usuario/bank-api/internal/bootstrap"
	customerApplication "github.com/seu-usuario/bank-api/internal/customer/application"
	customerDelivery "github.com/seu-usuario/bank-api/internal/customer/delivery"
	customerInfrastructure "github.com/seu-usuario/bank-api/internal/customer/infrastructure"
	"github.com/seu-usuario/bank-api/internal/database"
	sharedhttpmiddleware "github.com/seu-usuario/bank-api/internal/shared/http/middleware"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	bootstrap.Init()

	db := database.NewPool()
	log.Println("DB connected")

	// ======================
	// Config (fail-fast)
	// ======================
	appToken := os.Getenv("APP_TOKEN")
	if appToken == "" {
		log.Fatal("APP_TOKEN environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// ======================
	// Repositories
	// ======================
	customerRepo := customerInfrastructure.New(db)
	accountRepo := accountInfrastructure.New(db)

	userRepo := authInfrastructure.NewPostgresUserRepository(db)
	sessionRepo := authInfrastructure.NewPostgresSessionRepository(db)
	transactor := authInfrastructure.NewPostgresTransactor(db)

	// ======================
	// Services
	// ======================
	hasher := authInfrastructure.NewBcryptPasswordHasher(bcrypt.DefaultCost)
	tokenService := authInfrastructure.NewJWTTokenService(jwtSecret, 15*time.Minute)

	// ======================
	// Use Cases
	// ======================
	createAccountUC := accountApplication.NewCreateAccount(accountRepo, customerRepo, userRepo)
	depositUC := accountApplication.NewDeposit(accountRepo)
	withdrawUC := accountApplication.NewWithdraw(accountRepo)
	transferUC := accountApplication.NewTransfer(accountRepo)
	statementUC := accountApplication.NewGetStatement(accountRepo)

	registerUserUC := authApplication.NewRegisterUserUseCase(userRepo, customerRepo, hasher, transactor)
	loginUserUC := authApplication.NewLoginUserUseCase(userRepo, hasher, tokenService, sessionRepo)
	refreshAccessTokenUC := authApplication.NewRefreshAccessTokenUseCase(userRepo, tokenService, sessionRepo, transactor)
	getCurrentUserUC := authApplication.NewGetCurrentUserUseCase(userRepo)
	approveUserUC := authApplication.NewApproveUserUseCase(userRepo, accountRepo, customerRepo, transactor)

	getCustomerMeUC := customerApplication.NewGetCustomerMe(customerRepo)

	// ======================
	// Handlers
	// ======================
	accountHandler := accountDelivery.New(createAccountUC, depositUC, withdrawUC, transferUC, statementUC)
	authHandler := authDelivery.New(registerUserUC, loginUserUC, getCurrentUserUC, refreshAccessTokenUC, approveUserUC)
	customerHandler := customerDelivery.New(nil, getCustomerMeUC)

	// ======================
	// Middlewares
	// ======================
	appTokenMiddleware := sharedhttpmiddleware.AppToken(appToken)
	authMiddleware := authDelivery.NewJWTMiddleware(tokenService)

	withAuth := authMiddleware.RequireAuth

	// ======================
	// Routers
	// ======================

	// --- Auth Router ---
	authRouter := http.NewServeMux()

	// Onboarding (AppToken)
	authRouter.Handle("POST /auth/register", appTokenMiddleware(http.HandlerFunc(authHandler.Register)))
	authRouter.Handle("POST /auth/login", appTokenMiddleware(http.HandlerFunc(authHandler.Login)))

	// Authenticated (JWT)
	authRouter.Handle("POST /auth/refresh", withAuth(http.HandlerFunc(authHandler.Refresh)))
	authRouter.Handle("GET /auth/me", withAuth(http.HandlerFunc(authHandler.Me)))

	// --- API Router ---
	apiRouter := http.NewServeMux()
	apiRouter.Handle("POST /admin/users/{id}/approve", withAuth(http.HandlerFunc(authHandler.ApproveUser)))

	apiRouter.Handle("GET /customers/me", withAuth(http.HandlerFunc(customerHandler.Me)))

	apiRouter.Handle("POST /accounts", withAuth(http.HandlerFunc(accountHandler.CreateAccount)))
	apiRouter.Handle("POST /accounts/{id}/deposit", withAuth(http.HandlerFunc(accountHandler.Deposit)))
	apiRouter.Handle("POST /accounts/{id}/withdraw", withAuth(http.HandlerFunc(accountHandler.Withdraw)))
	apiRouter.Handle("GET /accounts/{id}/statement", withAuth(http.HandlerFunc(accountHandler.Statement)))
	apiRouter.Handle("POST /accounts/transfer", withAuth(http.HandlerFunc(accountHandler.Transfer)))

	// ======================
	// Main Router
	// ======================
	mainRouter := http.NewServeMux()

	mainRouter.Handle("/auth/", authRouter)
	mainRouter.Handle("/", apiRouter)

	log.Println("Server running in localhost on port 8080")

	if err := http.ListenAndServe(":8080", mainRouter); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
