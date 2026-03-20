package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type ctxKey struct{}

type Transaction struct {
	Tx *gorm.DB
}

type TransactionManager interface {
	Wrap(ctx context.Context, fn func(context.Context) error) error
}

type transactionManager struct {
	db *gorm.DB
}

func (tm *transactionManager) Wrap(ctx context.Context, fn func(context.Context) error) error {
	if _, ok := ctx.Value(ctxKey{}).(*Transaction); ok {
		return fn(ctx)
	}

	tx := tm.db.Begin()
	defer tx.Rollback()

	ctxWithTx := context.WithValue(ctx, ctxKey{}, &Transaction{Tx: tx})

	if err := fn(ctxWithTx); err != nil {
		return fmt.Errorf("fn execution: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &transactionManager{db: db}
}
