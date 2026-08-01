package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	query *Queries
	db    *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:    db,
		query: New(db),
	}
}

func (s *Store) execTx(context context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(context, nil)

	if err != nil {
		return err
	}

	q := New(tx)
	e := fn(q)
	if e != nil {
		if rbError := tx.Rollback(); rbError != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", e, rbError)
		}
		return e
	}
	return tx.Commit()
}

// Contains the input params of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// the result of transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		//record the transfer
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		//Create record entry from account
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		//Create record entry to account
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//TODO: update accounts' balance

		return nil
	})

	return result, err
}
