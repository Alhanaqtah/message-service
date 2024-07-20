package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (s *Storage) SaveMessage(ctx context.Context, msg *models.Message) (*models.Message, error) {
	const op = "storage.postgres.SaveMessage"

	row := s.pool.QueryRow(ctx, "INSERT INTO messages (content) VALUES ($1) RETURNING id, content", msg.Content)

	var message models.Message
	err := row.Scan(&message.ID, &message.Content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &message, nil
}

func (s *Storage) MarkMessageAsProcessed(ctx context.Context, id string) error {
	const op = "storage.postgres.MarkMessageAsProcessed"

	_, err := s.pool.Exec(ctx, "UPDATE messages SET status = 'processed', processed_at = $1 WHERE id = $2", time.Now(), id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) FetchStats(ctx context.Context) (*models.Stats, error) {
	const op = "storage.postgres.FetchStats"

	var stats models.Stats

	query := `
	SELECT 
	    COUNT(*) AS total_messages,
	    SUM(CASE WHEN status = 'processed' THEN 1 ELSE 0 END) AS processed_messages,
	    AVG(CASE WHEN status = 'processed' THEN EXTRACT(EPOCH FROM (processed_at - created_at)) ELSE NULL END) AS average_processing_time
	FROM 
	    messages;
	`

	var total sql.NullInt64
	var processed sql.NullInt64
	var avg sql.NullFloat64
	err := s.pool.QueryRow(ctx, query).Scan(&total, &processed, &avg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if total.Valid {
		stats.TotalMessages = int(total.Int64)
	} else {
		stats.TotalMessages = 0
	}

	if processed.Valid {
		stats.ProcessedMessages = int(processed.Int64)
	} else {
		stats.ProcessedMessages = 0
	}

	if avg.Valid {
		stats.AverageProcessingTime = time.Duration(avg.Float64 * 1000)
	} else {
		stats.AverageProcessingTime = 0
	}

	return &stats, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}
