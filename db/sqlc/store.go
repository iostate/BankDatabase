package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	// When you use the sql.DB database handle, you’re connecting with a built-in connection pool that creates and disposes of connections according
	// to your code’s needs. A handle through sql.DB is the most common way to do database access with Go
	// https://go.dev/doc/database/open-handle
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// #6 Takes context and callback fn as input
// Start new db transactopn
// Create new queries object and call the callbacl fn with the created Queries
// Finally commit or rollback the function based on the error from that function
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// if no isolation level is assigned in txoptions, read-committed is assigned which is default

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Call new fn wtith created transaction
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id`
	ToAccountID   int64 `json:"to_account_id`
	Amount        int64 `json:"amount`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account`
	ToAccount   Account  `json:"to_account`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Since we are accessing result from an inside fn
		// This becomes a closure which are often used when we want to get the result from
		// a callback fn
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Create 2 entries

		// FromAccount
		// Refer to models
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

// TransferTx will perform a xfer
// create xfer record, add acct. entries, upodate account balance within a single transaction
