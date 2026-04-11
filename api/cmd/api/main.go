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
	"golang.org/x/crypto/bcrypt"
)

func main() {
	bootstrap.Init()

	db := database.NewPool()

	log.Println("DB connected")

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
	hasher := authInfrastructure.NewBcryptPasswordHasher(bcrypt.DefaultCost)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-change-me"
	}
	tokenService := authInfrastructure.NewJWTTokenService(jwtSecret, 15*time.Minute)

	registerUserUC := authApplication.NewRegisterUserUseCase(userRepo, customerRepo, hasher)
	loginUserUC := authApplication.NewLoginUserUseCase(userRepo, hasher, tokenService, sessionRepo)
	refreshAccessTokenUC := authApplication.NewRefreshAccessTokenUseCase(userRepo, tokenService, sessionRepo)
	getCurrentUserUC := authApplication.NewGetCurrentUserUseCase(userRepo)
	getCustomerMeUC := customerApplication.NewGetCustomerMe(customerRepo)
	authHandler := authDelivery.New(registerUserUC, loginUserUC, getCurrentUserUC)
	authHandler.SetRefreshAccessTokenUseCase(refreshAccessTokenUC)
	customerHandler := customerDelivery.New(nil, getCustomerMeUC)
	authMiddleware := authDelivery.NewJWTMiddleware(tokenService)

	http.HandleFunc("POST /auth/register", authHandler.Register)
	http.HandleFunc("POST /auth/login", authHandler.Login)
	http.HandleFunc("POST /auth/refresh", authHandler.Refresh)
	http.Handle("GET /auth/me", authMiddleware.RequireAuth(http.HandlerFunc(authHandler.Me)))
	http.Handle("GET /customers/me", authMiddleware.RequireAuth(http.HandlerFunc(customerHandler.Me)))

	http.Handle("POST /accounts", authMiddleware.RequireAuth(http.HandlerFunc(accountHandler.CreateAccount)))
	http.Handle("POST /accounts/{id}/deposit", authMiddleware.RequireAuth(http.HandlerFunc(accountHandler.Deposit)))
	http.Handle("POST /accounts/{id}/withdraw", authMiddleware.RequireAuth(http.HandlerFunc(accountHandler.Withdraw)))
	http.Handle("GET /accounts/{id}/statement", authMiddleware.RequireAuth(http.HandlerFunc(accountHandler.Statement)))
	http.Handle("POST /accounts/transfer", authMiddleware.RequireAuth(http.HandlerFunc(accountHandler.Transfer)))

	log.Println("Server running on port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
