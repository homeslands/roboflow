package sqldb

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	_ SQLDB = &pgx.Conn{}
	_ SQLDB = pgx.Tx(nil)
	_ SQLDB = &pgxpool.Conn{}
	_ SQLDB = &pgxpool.Pool{}
	_ SQLDB = &pgxpool.Tx{}
)

// SQLDB is the common interface between *[pgx.Conn], *[pgx.Tx], *[pgxpool.Conn], *[pgxpool.Pool] and *[pgxpool.Tx].
type SQLDB interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row

	CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error)
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
}
