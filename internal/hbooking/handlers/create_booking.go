package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/spatecon/hbooking/internal/hbooking/domain"
)

type Booking struct {
	ID             string `json:"id,omitempty"`
	WorkshopID     int64  `json:"workshop_id"`
	ClientID       string `json:"client_id"`
	BeginAt        string `json:"begin_at"`
	EndAt          string `json:"end_at"`
	ClientTimezone string `json:"client_timezone"`
}

type CreateBookingRequest struct {
	Booking
}

func (h *Handlers) CreateBooking(c *gin.Context) {
	var req CreateBookingRequest
	err := c.ShouldBindJSON(&req)
	if BadRequest(c, err, "failed to bind request") {
		return
	}

	workshopID, err := strconv.ParseInt(c.Param("workshop_id"), 10, 64)
	if BadRequest(c, err, "failed to parse workshop_id") {
		return
	}

	booking := &domain.Booking{
		WorkshopID: workshopID,
		ClientID:   req.ClientID,
	}

	booking.ClientTimezone, err = time.LoadLocation(req.ClientTimezone)
	if BadRequest(c, err, "failed to load client timezone") {
		return
	}

	beginAt, err := time.Parse("02-01-2006 15:04", req.BeginAt)
	if BadRequest(c, err, "failed to parse begin_at") {
		return
	}

	endAt, err := time.Parse("02-01-2006 15:04", req.EndAt)
	if BadRequest(c, err, "failed to parse end_at") {
		return
	}

	booking.BeginAt = time.Date(
		beginAt.Year(),
		beginAt.Month(),
		beginAt.Day(),
		beginAt.Hour(), beginAt.Minute(), 0, 0,
		booking.ClientTimezone,
	).In(booking.ClientTimezone)

	booking.EndAt = time.Date(
		endAt.Year(),
		endAt.Month(),
		endAt.Day(),
		endAt.Hour(), endAt.Minute(), 0, 0,
		booking.ClientTimezone,
	).In(booking.ClientTimezone)

	now := time.Now().In(booking.ClientTimezone)
	if ValidationFailed(c, booking.BeginAt.Before(now), "begin_at is in the past") {
		return
	}

	if ValidationFailed(c, booking.EndAt.Before(booking.BeginAt), "end_at is before begin_at") {
		return
	}

	duration := booking.EndAt.Sub(booking.BeginAt)
	if ValidationFailed(c,
		duration < minBookingDuration || duration > maxBookingDuration,
		"invalid booking duration: must be in range [30m, 4h]",
	) {
		return
	}

	booking, err = h.repo.CreateBooking(c.Request.Context(), booking)
	if err != nil {
		if errors.Is(err, domain.ErrBookingOverlap) {
			Error(c, err, http.StatusBadRequest)
			return
		}
		if errors.Is(err, domain.ErrBookingOutOfWorkshopSchedule) {
			Error(c, err, http.StatusBadRequest)
			return
		}

		Error(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(200, Booking{
		ID:             booking.ID.String(),
		WorkshopID:     booking.WorkshopID,
		ClientID:       booking.ClientID,
		BeginAt:        booking.BeginAt.Format("02-01-2006 15:04"),
		EndAt:          booking.EndAt.Format("02-01-2006 15:04"),
		ClientTimezone: booking.ClientTimezone.String(),
	})
}
