package main

import (
	"log"
	"net/http"
	"time"

	accountApplication "github.com/seu-usuario/bank-api/internal/account/bankaccount/application"
	accountDelivery "github.com/seu-usuario/bank-api/internal/account/bankaccount/delivery"
	accountInfrastructure "github.com/seu-usuario/bank-api/internal/account/bankaccount/infrastructure"
	statementApplication "github.com/seu-usuario/bank-api/internal/account/statement/application"
	statementDelivery "github.com/seu-usuario/bank-api/internal/account/statement/delivery"
	statementInfrastructure "github.com/seu-usuario/bank-api/internal/account/statement/infrastructure"
	transactionApplication "github.com/seu-usuario/bank-api/internal/account/transaction/application"
	transactionDelivery "github.com/seu-usuario/bank-api/internal/account/transaction/delivery"
	transactionInfrastructure "github.com/seu-usuario/bank-api/internal/account/transaction/infrastructure"
	adminApplication "github.com/seu-usuario/bank-api/internal/admin/application"
	adminDelivery "github.com/seu-usuario/bank-api/internal/admin/delivery"
	authApplication "github.com/seu-usuario/bank-api/internal/auth/application"
	authDelivery "github.com/seu-usuario/bank-api/internal/auth/delivery"
	authInfrastructure "github.com/seu-usuario/bank-api/internal/auth/infrastructure"
	"github.com/seu-usuario/bank-api/internal/bootstrap"
	customerApplication "github.com/seu-usuario/bank-api/internal/customer/application"
	customerDelivery "github.com/seu-usuario/bank-api/internal/customer/delivery"
	customerInfrastructure "github.com/seu-usuario/bank-api/internal/customer/infrastructure"
	"github.com/seu-usuario/bank-api/internal/database"
	securityApplication "github.com/seu-usuario/bank-api/internal/security/application"
	securityDelivery "github.com/seu-usuario/bank-api/internal/security/delivery"
	securityDomain "github.com/seu-usuario/bank-api/internal/security/domain"
	securityInfrastructure "github.com/seu-usuario/bank-api/internal/security/infrastructure"
	sharedhttpmiddleware "github.com/seu-usuario/bank-api/internal/shared/http/middleware"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	bootstrap.Init()

	// ======================
	// Config (fail-fast)
	// ======================
	config := bootstrap.LoadConfig()

	db := database.NewPool()
	log.Println("DB connected")

	// ======================
	// Repositories
	// ======================
	customerRepo := customerInfrastructure.New(db)
	accountRepo := accountInfrastructure.New(db)
	transactionRepo := transactionInfrastructure.New(db)
	statementRepo := statementInfrastructure.New(db)

	userRepo := authInfrastructure.NewPostgresUserRepository(db)
	sessionRepo := authInfrastructure.NewPostgresSessionRepository(db)
	contactVerificationRepo := authInfrastructure.NewPostgresContactVerificationRepository(db)
	transactionPasswordRepo := securityInfrastructure.NewPostgresTransactionPasswordRepository(db)
	stepUpTokenRepo := securityInfrastructure.NewPostgresStepUpTokenRepository(db)
	transactor := authInfrastructure.NewPostgresTransactor(db)

	// ======================
	// Services
	// ======================
	hasher := authInfrastructure.NewBcryptPasswordHasher(bcrypt.DefaultCost)
	transactionPasswordHasher := securityInfrastructure.NewBcryptTransactionPasswordHasher(bcrypt.DefaultCost)
	tokenService := authInfrastructure.NewJWTTokenService(config.JWTSecret, 15*time.Minute)
	stepUpTokenSigner := securityInfrastructure.NewJWTStepUpTokenSigner(config.JWTSecret)

	// ======================
	// Use Cases
	// ======================
	branchPolicy := accountApplication.NewDefaultBranchPolicy()

	listAccountsUC := accountApplication.NewListAccounts(accountRepo)
	createAccountUC := accountApplication.NewCreateAccount(accountRepo, customerRepo, userRepo, branchPolicy)
	lookupInternalTransferRecipientsUC := accountApplication.NewLookupInternalTransferRecipients(accountRepo)
	depositUC := transactionApplication.NewDeposit(transactionRepo)
	withdrawUC := transactionApplication.NewWithdraw(transactionRepo)
	transferUC := transactionApplication.NewTransfer(transactionRepo)
	transferReceiptUC := transactionApplication.NewGetTransferReceipt(transactionRepo)
	statementUC := statementApplication.NewGetStatement(statementRepo)
	balanceUC := accountApplication.NewGetAccountBalance(accountRepo)

	registerUserUC := authApplication.NewRegisterUserUseCase(
		userRepo,
		customerRepo,
		customerRepo,
		contactVerificationRepo,
		hasher,
		transactor,
	)
	loginUserUC := authApplication.NewLoginUserUseCase(userRepo, accountRepo, hasher, tokenService, sessionRepo)
	refreshAccessTokenUC := authApplication.NewRefreshAccessTokenUseCase(userRepo, tokenService, sessionRepo, transactor)
	getCurrentUserUC := authApplication.NewGetCurrentUserUseCase(userRepo)
	requestContactVerificationUC := authApplication.NewRequestContactVerificationUseCase(contactVerificationRepo, userRepo)
	confirmContactVerificationUC := authApplication.NewConfirmContactVerificationUseCase(contactVerificationRepo)
	createTransactionPasswordUC := securityApplication.NewCreateTransactionPasswordUseCase(transactionPasswordRepo, userRepo, transactionPasswordHasher)
	authorizeStepUpUC := securityApplication.NewAuthorizeStepUpUseCase(
		transactionPasswordRepo,
		userRepo,
		transactionPasswordHasher,
		stepUpTokenRepo,
		stepUpTokenSigner,
		securityDomain.NewDefaultStepUpEndpointPolicy(),
	)
	approveUserUC := adminApplication.NewApproveUserUseCase(userRepo, accountRepo, customerRepo, transactor, branchPolicy)

	getCustomerMeUC := customerApplication.NewGetCustomerMe(customerRepo)
	checkCPFUC := customerApplication.NewCheckCPFUseCase(customerRepo)

	// ======================
	// Handlers
	// ======================
	accountHandler := accountDelivery.New(listAccountsUC, createAccountUC, balanceUC, lookupInternalTransferRecipientsUC)
	statementHandler := statementDelivery.New(statementUC)
	transactionHandler := transactionDelivery.New(depositUC, withdrawUC, transferUC, transferReceiptUC)
	authHandler := authDelivery.New(
		registerUserUC,
		loginUserUC,
		getCurrentUserUC,
		refreshAccessTokenUC,
		requestContactVerificationUC,
		confirmContactVerificationUC,
	)
	adminHandler := adminDelivery.New(approveUserUC)
	customerHandler := customerDelivery.New(nil, getCustomerMeUC, checkCPFUC)
	securityHandler := securityDelivery.New(createTransactionPasswordUC, authorizeStepUpUC)

	// ======================
	// Middlewares
	// ======================
	appTokenMiddleware := sharedhttpmiddleware.AppToken(config.AppToken)
	authMiddleware := authDelivery.NewJWTMiddleware(tokenService)

	withAuth := authMiddleware.RequireAuth

	// ======================
	// Routers
	// ======================

	// --- Auth Router ---
	authRouter := newAuthRouter(authHandler, customerHandler, appTokenMiddleware, withAuth)

	// --- API Router ---
	apiRouter := newAPIRouter(withAuth, adminHandler, accountHandler, customerHandler, statementHandler, transactionHandler, securityHandler)

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

func newAuthRouter(
	authHandler *authDelivery.Handler,
	customerHandler *customerDelivery.Handler,
	appTokenMiddleware func(http.Handler) http.Handler,
	withAuth func(http.Handler) http.Handler,
) *http.ServeMux {
	authRouter := http.NewServeMux()

	// Onboarding (AppToken)
	authRouter.Handle("POST /auth/cpf-check", appTokenMiddleware(http.HandlerFunc(customerHandler.CheckCPF)))
	authRouter.Handle("POST /auth/contact-verifications", appTokenMiddleware(http.HandlerFunc(authHandler.RequestContactVerification)))
	authRouter.Handle("POST /auth/contact-verifications/confirm", appTokenMiddleware(http.HandlerFunc(authHandler.ConfirmContactVerification)))
	authRouter.Handle("POST /auth/register", appTokenMiddleware(http.HandlerFunc(authHandler.Register)))
	authRouter.Handle("POST /auth/login", appTokenMiddleware(http.HandlerFunc(authHandler.Login)))

	// Session refresh is authenticated by the refresh token payload itself.
	authRouter.Handle("POST /auth/refresh", http.HandlerFunc(authHandler.Refresh))
	authRouter.Handle("GET /auth/me", withAuth(http.HandlerFunc(authHandler.Me)))

	return authRouter
}

func newAPIRouter(
	withAuth func(http.Handler) http.Handler,
	adminHandler *adminDelivery.Handler,
	accountHandler *accountDelivery.Handler,
	customerHandler *customerDelivery.Handler,
	statementHandler *statementDelivery.Handler,
	transactionHandler *transactionDelivery.Handler,
	securityHandler *securityDelivery.Handler,
) *http.ServeMux {
	apiRouter := http.NewServeMux()
	apiRouter.Handle("POST /admin/users/{id}/approve", withAuth(http.HandlerFunc(adminHandler.ApproveUser)))
	apiRouter.Handle("POST /admin/customers/{customer_id}/accounts", withAuth(http.HandlerFunc(accountHandler.CreateAccountForCustomer)))

	apiRouter.Handle("GET /customers/me", withAuth(http.HandlerFunc(customerHandler.Me)))

	apiRouter.Handle("GET /accounts", withAuth(http.HandlerFunc(accountHandler.ListAccounts)))
	apiRouter.Handle("GET /accounts/internal-transfers/recipients", withAuth(http.HandlerFunc(accountHandler.LookupInternalTransferRecipients)))
	apiRouter.Handle("GET /accounts/{id}/statement", withAuth(http.HandlerFunc(statementHandler.Statement)))
	apiRouter.Handle("GET /accounts/{id}/balance", withAuth(http.HandlerFunc(accountHandler.GetBalance)))
	apiRouter.Handle("POST /accounts/internal-transfers", withAuth(http.HandlerFunc(transactionHandler.Transfer)))
	apiRouter.Handle("GET /accounts/transfer/{transaction_reference}/receipt", withAuth(http.HandlerFunc(transactionHandler.TransferReceipt)))
	apiRouter.Handle("POST /security/transaction-password", withAuth(http.HandlerFunc(securityHandler.CreateTransactionPassword)))
	apiRouter.Handle("POST /security/step-up/authorize", withAuth(http.HandlerFunc(securityHandler.AuthorizeStepUp)))

	// Terminal cash operations are intentionally disabled until a real terminal channel exists.
	// apiRouter.Handle("POST /terminal/accounts/{id}/deposit", withAuth(http.HandlerFunc(transactionHandler.Deposit)))
	// apiRouter.Handle("POST /terminal/accounts/{id}/withdraw", withAuth(http.HandlerFunc(transactionHandler.Withdraw)))

	return apiRouter
}
