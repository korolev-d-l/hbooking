package handlers

import (
	"context"
	"log/slog"
	"time"

	"github.com/spatecon/hbooking/internal/hbooking/domain"
)

const (
	minBookingDuration = 30 * time.Minute
	maxBookingDuration = 4 * time.Hour
)

type Repository interface {
	CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	ListBookings(ctx context.Context, workshopID int64) ([]*domain.Booking, error)
}

type Handlers struct {
	logger *slog.Logger
	repo   Repository
}

func NewHandlers(logger *slog.Logger, repo Repository) *Handlers {
	return &Handlers{
		logger: logger,
		repo:   repo,
	}
}
