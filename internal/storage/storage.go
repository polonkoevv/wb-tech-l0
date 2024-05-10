package storage

import (
	"context"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/polonkoevv/wb-tech/internal/config"
)

func OpenDB(ctx context.Context, cfg config.Storage) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database,
	)
	log.Println(connectionString)
	db, err := sqlx.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("OpenDB sqlx open: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("OpenDB sqlx ping: %w", err)
	}

	return db, nil
}
