package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	DB          *pgxpool.Pool
	StmtBuilder sq.StatementBuilderType
}

func NewStore(p *pgxpool.Pool) *Store {
	return &Store{
		DB:          p,
		Queries:     New(p),
		StmtBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
