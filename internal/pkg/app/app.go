package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/spatecon/hbooking/internal/hbooking/handlers"
	"github.com/spatecon/hbooking/internal/hbooking/repository"
)

type App struct {
	cfg *config

	logger *slog.Logger

	dbConn     *pgxpool.Pool
	repository *repository.Repository

	router gin.IRouter
	http   *http.Server

	closers  []func() error
	handlers *handlers.Handlers
	closeCh  chan os.Signal
}

func NewApp() (*App, error) {
	app := new(App)
	app.initConfig()
	app.initLogger()

	err := app.initDBConn()
	if err != nil {
		return nil, fmt.Errorf("failed to init db connection: %w", err)
	}

	err = app.initServer()
	if err != nil {
		return nil, fmt.Errorf("failed to init http server: %w", err)
	}

	err = app.initGracefulShutdown()
	if err != nil {
		return nil, fmt.Errorf("failed to init graceful shutdown: %w", err)
	}

	return app, nil
}

func (a *App) Run() error {
	go func() {
		a.listenServer()
		a.stopWaitCloseSignal()
	}()

	a.waitCloseSignal()
	a.shutdown()

	return nil
}
