package db

import (
	"context"
	"errors"
)

// WithTx runs the given function in a transaction.
// If the function returns an error, the transaction is rolled back.
func (s *Store) WithTx(ctx context.Context, fn func(s Store) error) error {
	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}

	q := NewStore(s.DB)
	err = fn(*q)

	if err == nil {
		return tx.Commit(ctx)
	}

	if rbErr := tx.Rollback(ctx); rbErr != nil {
		return errors.Join(err, rbErr)
	}

	return err
}
