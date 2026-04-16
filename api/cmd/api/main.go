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

	appToken := os.Getenv("APP_TOKEN")
	if appToken == "" {
		log.Fatal("APP_TOKEN environment variable is required")
	}

	customerRepo := customerInfrastructure.New(db)

	accountRepo := accountInfrastructure.New(db)
	createAccountUC := accountApplication.NewCreateAccount(accountRepo, customerRepo)
	depositUC := accountApplication.NewDeposit(accountRepo)
	withdrawUC := accountApplication.NewWithdraw(accountRepo)
	transferUC := accountApplication.NewTransfer(accountRepo)
	statementUC := accountApplication.NewGetStatement(accountRepo)
	accountHandler := accountDelivery.New(createAccountUC, depositUC, withdrawUC, transferUC, statementUC)

	userRepo := authInfrastructure.NewPostgresUserRepository(db)
	sessionRepo := authInfrastructure.NewPostgresSessionRepository(db)
	transactor := authInfrastructure.NewPostgresTransactor(db)
	hasher := authInfrastructure.NewBcryptPasswordHasher(bcrypt.DefaultCost)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	tokenService := authInfrastructure.NewJWTTokenService(jwtSecret, 15*time.Minute)

	registerUserUC := authApplication.NewRegisterUserUseCase(userRepo, customerRepo, hasher)
	loginUserUC := authApplication.NewLoginUserUseCase(userRepo, hasher, tokenService, sessionRepo)
	refreshAccessTokenUC := authApplication.NewRefreshAccessTokenUseCase(userRepo, tokenService, sessionRepo, transactor)
	getCurrentUserUC := authApplication.NewGetCurrentUserUseCase(userRepo)
	getCustomerMeUC := customerApplication.NewGetCustomerMe(customerRepo)
	authHandler := authDelivery.New(registerUserUC, loginUserUC, getCurrentUserUC, refreshAccessTokenUC)
	customerHandler := customerDelivery.New(nil, getCustomerMeUC)
	authMiddleware := authDelivery.NewJWTMiddleware(tokenService)

	router := http.NewServeMux()

	router.HandleFunc("POST /auth/register", authHandler.Register)
	router.HandleFunc("POST /auth/login", authHandler.Login)
	router.HandleFunc("POST /auth/refresh", authHandler.Refresh)

	router.Handle("GET /auth/me", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.Me)))
	router.Handle("GET /customers/me", authMiddleware.RequireAuth(http.HandlerFunc(customerHandler.Me)))

	router.Handle("POST /accounts", authMiddleware.RequireAuth(http.HandlerFunc(accountHandler.CreateAccount)))
	router.Handle("POST /accounts/{id}/deposit", authMiddleware.RequireAuth(http.HandlerFunc(accountHandler.Deposit)))
	router.Handle("POST /accounts/{id}/withdraw", authMiddleware.RequireAuth(http.HandlerFunc(accountHandler.Withdraw)))
	router.Handle("GET /accounts/{id}/statement", authMiddleware.RequireAuth(http.HandlerFunc(accountHandler.Statement)))
	router.Handle("POST /accounts/transfer", authMiddleware.RequireAuth(http.HandlerFunc(accountHandler.Transfer)))

	handler := http.Handler(router)
	handler = sharedhttpmiddleware.AppToken(appToken)(handler)

	log.Println("Server running in localhost on port 8080")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
