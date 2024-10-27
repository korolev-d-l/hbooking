package repository

import (
	"log/slog"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	logger *slog.Logger
	pool   *pgxpool.Pool
}

func NewRepository(logger *slog.Logger, pool *pgxpool.Pool) *Repository {
	return &Repository{
		logger: logger,
		pool:   pool,
	}
}
