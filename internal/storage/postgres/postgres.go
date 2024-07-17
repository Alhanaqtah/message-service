package postgres

import (
	"message-service/internal/config/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(cfg *config.Storage) *Storage {
	return &Storage{}
}
