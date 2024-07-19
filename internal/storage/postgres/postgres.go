package postgres

import (
	"context"
	"fmt"

	"message-service/internal/config/config"
	"message-service/internal/models"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(cfg *config.Storage) (*Storage, error) {
	const op = "storage.postgres.New"

	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db := stdlib.OpenDB(*pool.Config().ConnConfig)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) SaveMessage(ctx context.Context, msg *models.Message) error {
	const op = "storage.postgres.SaveMessage"

	_, err := s.pool.Exec(ctx, "INSERT INTO messages (content) VALUES ($1)", msg.Content)
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Close() {
	s.pool.Close()
}
