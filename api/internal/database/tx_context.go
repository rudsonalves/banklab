package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// This file defines utilities for managing database transactions within a context.
// It allows passing a pgx.Tx through the context, enabling transaction-aware
// operations without explicitly threading the transaction through function parameters.
//
// Example usage:
//
//	func SomeUseCase(ctx context.Context) error {
//	    tx, err := db.Begin(ctx)
//	    if err != nil {
//	        return err
//	    }
//	    defer tx.Rollback(ctx)
//
//	    ctx = database.ContextWithTx(ctx, tx)
//
//	    // Perform database operations using TxFromContext to get the transaction
//	    // ...
//
//	    return tx.Commit(ctx)
//	}
type txContextKey struct{}

// ContextWithTx returns a new context derived from ctx that carries the given pgx.Tx.
func ContextWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}

// TxFromContext retrieves the pgx.Tx from the context.
// It returns the transaction and true if found, otherwise nil and false.
func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txContextKey{}).(pgx.Tx)
	return tx, ok
}
