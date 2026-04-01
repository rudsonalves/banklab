package main

import (
	"log"
	"net/http"

	"github.com/seu-usuario/bank-api/internal/customer/handler"
	"github.com/seu-usuario/bank-api/internal/customer/repository"
	"github.com/seu-usuario/bank-api/internal/customer/usecase"
	"github.com/seu-usuario/bank-api/internal/database"
)

func main() {
	db := database.NewPool()

	log.Println("DB connected")

	repo := repository.New(db)
	uc := usecase.NewCreateCustomer(repo)
	handler := handler.New(uc)

	http.HandleFunc("/customers", handler.Create)

	log.Println("Server running on port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
