package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (a *App) initDBConn() error {
	cfg, err := pgxpool.ParseConfig(a.cfg.DB.URL)
	if err != nil {
		return err
	}

	db, err := pgxpool.NewWithConfig(context.TODO(), cfg)
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
