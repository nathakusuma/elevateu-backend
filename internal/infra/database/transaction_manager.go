package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ITransactionManager interface {
	BeginTx(ctx context.Context) (ITransaction, error)
}

type transactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) ITransactionManager {
	return &transactionManager{
		db: db,
	}
}

func (t *transactionManager) BeginTx(ctx context.Context) (ITransaction, error) {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return newTransactionWrapper(tx), nil
}
