package postgres

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
)

type DB struct {
	db *sql.DB
}

func NewRepository(databaseDSN string) (*DB, error) {
	if databaseDSN == "" {
		log.Println("No database DSN provided")
		return &DB{db: nil}, nil
	}

	db, err := sql.Open("pgx", databaseDSN)

	if err != nil {
		return nil, err
	}

	repo := &DB{db: db}

	return repo, nil
}

func (repo *DB) Ping(ctx context.Context) error {
	if repo.db == nil {
		return errors.New("db is nil")
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return repo.db.PingContext(ctx)
}

func (repo *DB) Close() {
	if repo.db == nil {
		return
	}

	err := repo.db.Close()
	if err != nil {
		log.Printf("Failed to close database: %v", err)
	}
}
