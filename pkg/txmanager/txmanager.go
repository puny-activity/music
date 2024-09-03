package txmanager

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type TxManager struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *TxManager {
	return &TxManager{
		db: db,
	}
}

func (m *TxManager) Transaction(ctx context.Context, f func(context.Context, *sqlx.Tx) error) error {
	tx, errBegin := m.db.BeginTxx(ctx, nil)
	if errBegin != nil {
		return fmt.Errorf("failed to begin transaction: %w", errBegin)
	}

	errFunction := f(ctx, tx)
	if errFunction != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			return fmt.Errorf("failed to rollback transaction: %w on error: %w", errRollback, errFunction)
		}
		return errFunction
	}

	errCommit := tx.Commit()
	if errCommit != nil {
		return fmt.Errorf("failed to commit transaction: %w", errCommit)
	}

	return nil
}
