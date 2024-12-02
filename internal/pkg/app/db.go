package app

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (a *App) initDBConn() error {
	databaseURL := os.Getenv("DATABASE_URL")

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return err
	}

	config.MaxConns = 10
	config.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return err
	}
	a.closers = append(a.closers, func() error {
		db.Close()
		return nil
	})

	a.dbConn = db

	return nil
}
