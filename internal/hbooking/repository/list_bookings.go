package repository

import (
	"context"
	"fmt"

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
	defer rows.Close()

	bookings := make([]*domain.Booking, 0)
	for rows.Next() {
		var (
			b       Booking
			booking *domain.Booking
		)
		err = rows.Scan(
			&b.ID,
			&b.WorkshopID,
			&b.ClientID,
			&b.BeginAt,
			&b.EndAt,
			&b.ClientTimezone,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}

		booking, err = ConvertFromPGBooking(&b)
		if err != nil {
			return nil, fmt.Errorf("failed to convert booking: %w", err)
		}

		bookings = append(bookings, booking)
	}

	return bookings, nil
}
