package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"

	"github.com/spatecon/hbooking/internal/hbooking/domain"
)

type Booking struct {
	ID             uuid.UUID
	WorkshopID     int64
	ClientID       string
	BeginAt        pgtype.Timestamp
	EndAt          pgtype.Timestamp
	ClientTimezone string
}

func ConvertToPGBooking(b *domain.Booking) *Booking {
	return &Booking{
		ID:             b.ID,
		WorkshopID:     b.WorkshopID,
		ClientID:       b.ClientID,
		BeginAt:        pgtype.Timestamp{Time: b.BeginAt.UTC(), Status: pgtype.Present},
		EndAt:          pgtype.Timestamp{Time: b.EndAt.UTC(), Status: pgtype.Present},
		ClientTimezone: b.ClientTimezone.String(),
	}
}

func ConvertFromPGBooking(b *Booking) (*domain.Booking, error) {
	tz, err := time.LoadLocation(b.ClientTimezone)
	if err != nil {
		return nil, fmt.Errorf("failed to load client timezone: %w", err)
	}

	return &domain.Booking{
		ID:             b.ID,
		WorkshopID:     b.WorkshopID,
		ClientID:       b.ClientID,
		BeginAt:        b.BeginAt.Time.In(tz),
		EndAt:          b.EndAt.Time.In(tz),
		ClientTimezone: tz,
	}, nil
}

type WorkshopSchedule struct {
	WorkshopID int64
	Timezone   string
	BeginAt    pgtype.Time
	EndAt      pgtype.Time
}

func (r *Repository) CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	var err error

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		txErr := tx.Rollback(ctx)
		if txErr != nil && !errors.Is(txErr, pgx.ErrTxClosed) {
			r.logger.Warn("failed to rollback transaction", "error", txErr)
		}
	}()

	// Check if booking is out of workshop schedule
	query := `
		SELECT workshop_id, workshop_timezone, begin_at, end_at
		FROM workshop_schedules
		WHERE workshop_id = $1
`

	var schedule WorkshopSchedule
	err = tx.QueryRow(ctx, query, booking.WorkshopID).Scan(
		&schedule.WorkshopID,
		&schedule.Timezone,
		&schedule.BeginAt,
		&schedule.EndAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query workshop schedule: %w", err)
	}

	if !r.bookingComplyWithSchedule(booking, schedule) {
		return nil, domain.ErrBookingOutOfWorkshopSchedule
	}

	// Check if booking overlaps with another booking
	query = `
		SELECT COUNT(booking_id)
		FROM workshop_bookings
		WHERE
			workshop_id = $1 AND
			((begin_at, end_at) OVERLAPS ($2, $3));
`
	pgBooking := ConvertToPGBooking(booking)

	var count int
	err = tx.QueryRow(ctx, query,
		pgBooking.WorkshopID,
		pgBooking.BeginAt,
		pgBooking.EndAt,
	).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings: %w", err)
	}

	if count > 0 {
		return nil, domain.ErrBookingOverlap
	}

	// Insert booking into database
	query = `
		INSERT INTO workshop_bookings 
		    (workshop_id, begin_at, end_at, client_id, client_timezone)
		VALUES
		    ($1, $2, $3, $4, $5)
		RETURNING booking_id;
`

	err = tx.QueryRow(ctx, query,
		pgBooking.WorkshopID,
		pgBooking.BeginAt,
		pgBooking.EndAt,
		pgBooking.ClientID,
		pgBooking.ClientTimezone,
	).Scan(&pgBooking.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert booking: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	booking, err = ConvertFromPGBooking(pgBooking)
	if err != nil {
		return nil, fmt.Errorf("failed to convert booking: %w", err)
	}

	return booking, nil
}

func (r *Repository) bookingComplyWithSchedule(booking *domain.Booking, ws WorkshopSchedule) bool {
	wsTz, err := time.LoadLocation(ws.Timezone)
	if err != nil {
		r.logger.Error("failed to load workshop timezone", "error", err)
		return false
	}

	wsBeginAtMinutes := ws.BeginAt.Microseconds / 1e6 / 60
	wsEndAtMinutes := ws.EndAt.Microseconds / 1e6 / 60

	wsBeginAtBookingDate := time.Date(
		booking.BeginAt.Year(), booking.BeginAt.Month(), booking.BeginAt.Day(),
		int(wsBeginAtMinutes/60), int(wsBeginAtMinutes%60), 0, 0, wsTz,
	)
	wsEndAtBookingDate := time.Date(
		booking.EndAt.Year(), booking.EndAt.Month(), booking.EndAt.Day(),
		int(wsEndAtMinutes/60), int(wsEndAtMinutes%60), 0, 0, wsTz,
	)
	if wsEndAtBookingDate.Before(wsBeginAtBookingDate) { // next day
		wsEndAtBookingDate = wsEndAtBookingDate.AddDate(0, 0, 1)
	}

	// check booking begin_at and end_at are within workshop schedule
	if booking.BeginAt.Before(wsBeginAtBookingDate) || booking.EndAt.After(wsEndAtBookingDate) {
		return false
	}

	return true
}
