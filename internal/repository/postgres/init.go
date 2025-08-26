package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

func (repo *DB) Migrate() {
	if repo.db == nil {
		return
	}

	driver, err := pgmigrate.WithInstance(repo.db, &pgmigrate.Config{})
	if err != nil {
		log.Fatalf("could not create postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("could not start migration: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("migrations applied successfully")
}

func (repo *DB) Ping(ctx context.Context) error {
	if repo.db == nil {
		return errors.New("db is nil")
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return repo.db.PingContext(ctx)
}
