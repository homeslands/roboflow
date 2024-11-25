package db

import "github.com/jackc/pgx/v5/pgxpool"

type Store struct {
	p *pgxpool.Pool
	*Queries
}

func NewStore(p *pgxpool.Pool) *Store {
	return &Store{
		p:       p,
		Queries: New(p),
	}
}
