package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseStore defines all functions to execute sql queries and transactions.
type DatabaseStore interface {
	Querier
}

// PostgresDBStore holds the database connections pool, and provides all functions
// to execute postgrSQL queries and transactions.
type PostgresDBStore struct {
	Pool *pgxpool.Pool
	*Queries
}

const (
	MaxDbOpenConns    = 10
	MaxDbConnIdleTime = 5 * time.Second
	MaxDbConnLifetime = 5 * time.Minute
)

// NewPostgresDBStore creates a new DatabaseStore
func NewPostgresDBStore(connString string) (DatabaseStore, error) {
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = MaxDbOpenConns
	poolConfig.MaxConnIdleTime = MaxDbConnIdleTime
	poolConfig.MaxConnLifetime = MaxDbConnLifetime

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &PostgresDBStore{
		Pool:    db,
		Queries: New(db),
	}, nil
}
