package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrBookingOutOfWorkshopSchedule = errors.New("booking is out of workshop schedule")
	ErrBookingOverlap               = errors.New("booking overlaps with another booking")
)

type Booking struct {
	ID             uuid.UUID
	WorkshopID     int64
	ClientID       string
	BeginAt        time.Time
	EndAt          time.Time
	ClientTimezone *time.Location
}
