package handlers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/spatecon/hbooking/internal/hbooking/domain"
)

type CreateBookingResponse struct {
	ID             string `json:"id,omitempty"`
	WorkshopID     int64  `json:"workshop_id"`
	ClientID       string `json:"client_id"`
	BeginAt        string `json:"begin_at"`
	EndAt          string `json:"end_at"`
	ClientTimezone string `json:"client_timezone"`
}

type CreateBookingRequest struct {
	WorkshopID     int64  `json:"workshop_id"`
	ClientID       string `json:"client_id"`
	BeginAt        string `json:"begin_at"`
	EndAt          string `json:"end_at"`
	ClientTimezone string `json:"client_timezone"`
}

func (h *Handlers) CreateBooking(c *gin.Context) {
	var req CreateBookingRequest
	if BadRequest(c, makeCreateBookingRequest(c, &req)) {
		return
	}

	var booking domain.Booking
	if BadRequest(c, requestToDomainBooking(&req, &booking)) {
		return
	}

	b, err := createBookingService(c.Request.Context(), h.repo, &booking)
	if err != nil {
		if ValidationFailed(c, err) {
			return
		}

		InternalServerError(c, err)
		return
	}

	c.JSON(200, makeCreateBookingResponse(b))
}

func makeCreateBookingRequest(c *gin.Context, r *CreateBookingRequest) error {
	var err error
	if err = c.ShouldBindJSON(r); err != nil {
		return fmt.Errorf("failed to bind request: %w", err)
	}

	if r.WorkshopID, err = strconv.ParseInt(c.Param("workshop_id"), 10, 64); err != nil {
		return fmt.Errorf("failed to parse workshop_id")
	}

	return nil
}

func requestToDomainBooking(req *CreateBookingRequest, booking *domain.Booking) error {
	clientTimezone, err := time.LoadLocation(req.ClientTimezone)
	if err != nil {
		return fmt.Errorf("failed to load client timezone")
	}

	beginAt, err := time.Parse("02-01-2006 15:04", req.BeginAt)
	if err != nil {
		return fmt.Errorf("failed to parse begin_at")
	}

	endAt, err := time.Parse("02-01-2006 15:04", req.EndAt)
	if err != nil {
		return fmt.Errorf("failed to parse end_at")
	}

	booking.WorkshopID = req.WorkshopID
	booking.ClientID = req.ClientID
	booking.BeginAt = time.Date(
		beginAt.Year(),
		beginAt.Month(),
		beginAt.Day(),
		beginAt.Hour(), beginAt.Minute(), 0, 0,
		clientTimezone,
	).In(clientTimezone)
	booking.EndAt = time.Date(
		endAt.Year(),
		endAt.Month(),
		endAt.Day(),
		endAt.Hour(), endAt.Minute(), 0, 0,
		clientTimezone,
	).In(clientTimezone)
	booking.ClientTimezone = clientTimezone

	return nil
}

func createBookingService(ctx context.Context, repo Repository, booking *domain.Booking) (*domain.Booking, error) {
	now := time.Now().In(booking.ClientTimezone)
	if booking.BeginAt.Before(now) {
		return nil, ValidationErrorStr("begin_at is in the past")
	}

	if booking.EndAt.Before(booking.BeginAt) {
		return nil, ValidationErrorStr("end_at is before begin_at")
	}

	duration := booking.EndAt.Sub(booking.BeginAt)
	if duration < minBookingDuration || duration > maxBookingDuration {
		return nil, ValidationErrorStr("invalid booking duration: must be in range [30m, 4h]")
	}

	b, err := repo.CreateBooking(ctx, booking)
	if err != nil {
		if errors.Is(err, domain.ErrBookingOverlap) || errors.Is(err, domain.ErrBookingOutOfWorkshopSchedule) {
			return nil, ValidationError(err)
		}

		return nil, err
	}
	return b, nil
}

func makeCreateBookingResponse(b *domain.Booking) CreateBookingResponse {
	return CreateBookingResponse{
		ID:             b.ID.String(),
		WorkshopID:     b.WorkshopID,
		ClientID:       b.ClientID,
		BeginAt:        b.BeginAt.Format("02-01-2006 15:04"),
		EndAt:          b.EndAt.Format("02-01-2006 15:04"),
		ClientTimezone: b.ClientTimezone.String(),
	}
}
