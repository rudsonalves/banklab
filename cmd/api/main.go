package main

import (
	"log"
	"net/http"

	accountApplication "github.com/seu-usuario/bank-api/internal/account/application"
	accountDelivery "github.com/seu-usuario/bank-api/internal/account/delivery"
	accountInfrastructure "github.com/seu-usuario/bank-api/internal/account/infrastructure"
	customerApplication "github.com/seu-usuario/bank-api/internal/customer/application"
	customerDelivery "github.com/seu-usuario/bank-api/internal/customer/delivery"
	customerInfrastructure "github.com/seu-usuario/bank-api/internal/customer/infrastructure"
	"github.com/seu-usuario/bank-api/internal/database"
)

func main() {
	db := database.NewPool()

	log.Println("DB connected")

	customerRepo := customerInfrastructure.New(db)
	customerUC := customerApplication.NewCreateCustomer(customerRepo)
	customerHandler := customerDelivery.New(customerUC)

	accountRepo := accountInfrastructure.New(db)
	createAccountUC := accountApplication.NewCreateAccount(accountRepo, customerRepo)
	depositUC := accountApplication.NewDeposit(accountRepo)
	accountHandler := accountDelivery.New(createAccountUC, depositUC)

	http.HandleFunc("POST /customers", customerHandler.Create)
	http.HandleFunc("POST /accounts", accountHandler.CreateAccount)
	http.HandleFunc("POST /accounts/{id}/deposit", accountHandler.Deposit)

	log.Println("Server running on port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
