package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool() *pgxpool.Pool {
	connString := "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal("unable to connect to database:", err)
	}

	return pool
}
