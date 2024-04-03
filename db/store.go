package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseStore defines all functions to execute sql queries and transactions.
type DatabaseStore interface {
	Querier
	CreateReservationTx(ctx context.Context, arg CreateReservationParams) (Reservation, error)
}

// PostgresDBStore holds the database connections pool, and provides all functions
// to execute postgrSQL queries and transactions.
type PostgresDBStore struct {
	DBConnPool *pgxpool.Pool
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

	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &PostgresDBStore{
		DBConnPool: conn,
		Queries:    New(conn),
	}, nil
}

// execTx executes a function within a database transaction
func (store *PostgresDBStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.DBConnPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
