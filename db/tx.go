package db

import (
	"context"
	"errors"
)

func (s *Store) WithTx(ctx context.Context, fn func(s Store) error) error {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return err
	}

	q := NewStore(s.p)
	err = fn(*q)

	if err == nil {
		return tx.Commit(ctx)
	}

	if rbErr := tx.Rollback(ctx); rbErr != nil {
		return errors.Join(err, rbErr)
	}

	return err
}
