package database

import (
	"context"
	"log"
	"net"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// This file defines the database connection pool and transaction context utilities.
// It provides a function to create a new connection pool and utilities to manage
// transactions within a context, allowing for cleaner transaction handling in the application.
//
// The NewPool function establishes a connection to the PostgreSQL database using
// the pgxpool package and returns a connection pool that can be used throughout the application.
func NewPool(config Config) *pgxpool.Pool {
	connString := connectionString(config)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal("unable to connect to database:", err)
	}

	return pool
}

func connectionString(config Config) string {
	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.User, config.Password),
		Host:   net.JoinHostPort(config.Host, config.Port),
		Path:   config.Name,
	}

	query := databaseURL.Query()
	query.Set("sslmode", "disable")
	databaseURL.RawQuery = query.Encode()

	return databaseURL.String()
}
