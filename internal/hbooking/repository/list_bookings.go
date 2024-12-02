package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/spatecon/hbooking/internal/hbooking/domain"
)

func (r *Repository) ListBookings(ctx context.Context, workshopID int64) ([]*domain.Booking, error) {
	query := `	
		SELECT booking_id, workshop_id, client_id, begin_at, end_at, client_timezone
		FROM workshop_bookings
		WHERE workshop_id = $1;
	`

	rows, err := r.pool.Query(ctx, query, workshopID)
	if err != nil {
		return nil, err
	}

	bookings, err := pgx.CollectRows(rows, pgx.RowToStructByName[Booking])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows booking: %w", err)
	}

	bookingsDomain := make([]*domain.Booking, len(bookings))
	for i, b := range bookings {
		bookingsDomain[i], err = ConvertFromPGBooking(&b)
		if err != nil {
			return nil, fmt.Errorf("failed to convert booking: %w", err)
		}
	}

	return bookingsDomain, nil
}
