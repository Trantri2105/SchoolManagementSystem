package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log"
)

type TransactionManager interface {
	ExecTransaction(ctx context.Context, callback func(ctx context.Context, tx *sqlx.Tx) error) error
}

type transactionManager struct {
	db *sqlx.DB
}

func (t *transactionManager) ExecTransaction(ctx context.Context, callback func(ctx context.Context, tx *sqlx.Tx) error) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println("TransactionManager, begin transaction error:", err)
		return err
	}

	defer func() {
		if err != nil {
			log.Println("TransactionManager, exec transaction error:", err)
			if err = tx.Rollback(); err != nil {
				log.Println("TransactionManager, rollback transaction error:", err)
			}
		}
	}()

	err = callback(ctx, tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func NewTransactionManager(db *sqlx.DB) TransactionManager {
	return &transactionManager{db: db}
}
