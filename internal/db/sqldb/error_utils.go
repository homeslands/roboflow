package sqldb

import (
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// IsUniqueViolationError checks if the error is a PostgreSQL unique constraint violation.
func IsUniqueViolationError(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, constraint) {
		return true
	}
	return false
}

// IsNoRowsError checks if the error is a pgx.ErrNoRows.
func IsNoRowsError(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
