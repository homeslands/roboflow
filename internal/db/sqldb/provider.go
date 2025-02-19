package sqldb

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Provider = (*sqlDBProvider)(nil)

// Provider represents a provider of SQL database connections.
type Provider interface {
	// DB returns the SQLDB instance.
	DB() SQLDB

	// WithTx executes a function in a new transaction.
	WithTx(ctx context.Context, fn func(db SQLDB) error) error
}

type sqlDBProvider struct {
	pool *pgxpool.Pool
}

// NewProvider creates a new SQLDBProvider.
//
//nolint:revive
func NewProvider(pool *pgxpool.Pool) *sqlDBProvider {
	return &sqlDBProvider{pool}
}

func (p *sqlDBProvider) DB() SQLDB {
	return p.pool
}

func (p *sqlDBProvider) WithTx(ctx context.Context, txFunc func(SQLDB) error) (err error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			rbErr := tx.Rollback(ctx)
			if !errors.Is(rbErr, pgx.ErrTxClosed) {
				err = errors.Join(err, rbErr)
			}
		}
	}()

	if err = txFunc(tx); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		err = fmt.Errorf("commit transaction: %w", err)
	}

	return err
}
