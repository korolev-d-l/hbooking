package handlers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/spatecon/hbooking/internal/hbooking/domain"
)

type ListBookingsRequest struct {
	WorkshopID int64 `json:"workshop_id"`
}

type ListBookingsResponse struct {
	Bookings []ListBookingResponse `json:"bookings"`
}

type ListBookingResponse struct {
	WorkshopID     int64  `json:"workshop_id"`
	ClientID       string `json:"client_id"`
	BeginAt        string `json:"begin_at"`
	EndAt          string `json:"end_at"`
	ClientTimezone string `json:"client_timezone"`
}

func (h *Handlers) ListBookings(c *gin.Context) {
	var req ListBookingsRequest
	if BadRequest(c, makeListBookingsRequest(c, &req)) {
		return
	}

	bookings, err := h.repo.ListBookings(c.Request.Context(), req.WorkshopID)
	if InternalServerError(c, err) {
		return
	}

	c.JSON(200, makeListBookingsResponse(bookings))

}

func makeListBookingsRequest(c *gin.Context, r *ListBookingsRequest) error {
	var err error
	if r.WorkshopID, err = strconv.ParseInt(c.Param("workshop_id"), 10, 64); err != nil {
		return fmt.Errorf("failed to parse workshop_id")
	}

	return nil
}

func makeListBookingsResponse(bookings []*domain.Booking) ListBookingsResponse {
	respBookings := make([]ListBookingResponse, len(bookings))
	for i, b := range bookings {
		respBookings[i] = ListBookingResponse{
			WorkshopID:     b.WorkshopID,
			ClientID:       b.ClientID,
			BeginAt:        b.BeginAt.Format("02-01-2006 15:04"),
			EndAt:          b.EndAt.Format("02-01-2006 15:04"),
			ClientTimezone: b.ClientTimezone.String(),
		}
	}
	return ListBookingsResponse{Bookings: respBookings}
}
