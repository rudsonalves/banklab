package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// This file defines the database connection pool and transaction context utilities.
// It provides a function to create a new connection pool and utilities to manage
// transactions within a context, allowing for cleaner transaction handling in the application.
//
// The NewPool function establishes a connection to the PostgreSQL database using
// the pgxpool package and returns a connection pool that can be used throughout the application.
func NewPool() *pgxpool.Pool {
	connString := "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable"

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal("unable to connect to database:", err)
	}

	return pool
}
