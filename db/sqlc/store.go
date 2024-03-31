package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	Querier
	CreateUserTransferTx(ctx context.Context, arg CreateUserTransferTxParams) (CreateUserTransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{Queries: New(db),
		db: db,
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return fmt.Errorf("err: %v, txErr: %v", err, txErr)
		}
		return err
	}

	return tx.Commit()
}
