package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewConnection(dsn string) (*PostgresDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	log.Println("Conectado ao PostgreSQL com sucesso.")

	return &PostgresDB{Pool: pool}, nil
}

// Close fecha o pool de conexões
func (db *PostgresDB) Close() {
	db.Pool.Close()
}
