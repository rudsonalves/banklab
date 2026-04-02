package main

import (
	"log"
	"net/http"

	"github.com/seu-usuario/bank-api/internal/customer/application"
	"github.com/seu-usuario/bank-api/internal/customer/delivery"
	"github.com/seu-usuario/bank-api/internal/customer/infrastructure"
	"github.com/seu-usuario/bank-api/internal/database"
)

func main() {
	db := database.NewPool()

	log.Println("DB connected")

	repo := infrastructure.New(db)
	uc := application.NewCreateCustomer(repo)
	handler := delivery.New(uc)

	http.HandleFunc("/customers", handler.Create)

	log.Println("Server running on port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
