package database

import "github.com/jmoiron/sqlx"

type ITransaction interface {
	GetTx() *sqlx.Tx
	Commit() error
	Rollback() error
}

type transactionWrapper struct {
	tx *sqlx.Tx
}

func newTransactionWrapper(tx *sqlx.Tx) ITransaction {
	return &transactionWrapper{
		tx: tx,
	}
}

func (t *transactionWrapper) GetTx() *sqlx.Tx {
	return t.tx
}

func (t *transactionWrapper) Commit() error {
	return t.tx.Commit()
}

func (t *transactionWrapper) Rollback() error {
	return t.tx.Rollback()
}
