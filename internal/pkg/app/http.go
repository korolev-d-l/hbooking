package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/spatecon/hbooking/internal/hbooking/handlers"
	"github.com/spatecon/hbooking/internal/hbooking/repository"
)

func (a *App) initServer() error {
	a.repository = repository.NewRepository(a.logger, a.dbConn)

	err := a.repository.Migrate()
	if err != nil {
		return fmt.Errorf("failed to migrate db: %w", err)
	}

	a.handlers = handlers.NewHandlers(a.logger, a.repository)

	gin.SetMode(a.cfg.HTTP.GinMode)
	router := gin.Default()
	router.POST("/api/v1/bookings/:workshop_id", a.handlers.CreateBooking)
	router.GET("/api/v1/bookings/:workshop_id", a.handlers.ListBookings)

	a.router = router

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", a.cfg.HTTP.Server.Port),
		Handler:           http.TimeoutHandler(router, a.cfg.HTTP.Server.HandlerTimeout, "Timeout"),
		ErrorLog:          slog.NewLogLogger(a.logger.Handler(), slog.LevelError),
		ReadHeaderTimeout: a.cfg.HTTP.Server.ReadHeaderTimeout,
		ReadTimeout:       a.cfg.HTTP.Server.ReadTimeout,
		IdleTimeout:       a.cfg.HTTP.Server.IdleTimeout,
	}

	a.http = server
	a.closers = append(a.closers, func() error {
		ctx, cancel := context.WithTimeout(context.Background(), a.cfg.HTTP.Server.ShutdownTimeout)
		defer cancel()

		return a.http.Shutdown(ctx)
	})

	return nil
}

func (a *App) listenServer() {
	err := a.http.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.logger.Error("http server failed", "error", err)
	}
}
