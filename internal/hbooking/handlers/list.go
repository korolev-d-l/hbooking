package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ListBookingsRequest struct {
	WorkshopID int `json:"workshop_id"`
}

type ListBookingsResponse struct {
	Bookings []*Booking `json:"bookings"`
}

func (h *Handlers) ListBookings(c *gin.Context) {
	workshopID, err := strconv.ParseInt(c.Param("workshop_id"), 10, 64)
	if BadRequest(c, err, "failed to parse workshop_id") {
		return
	}

	bookings, err := h.repo.ListBookings(c.Request.Context(), workshopID)
	if Error(c, err, http.StatusInternalServerError) {
		return
	}

	respBookings := make([]*Booking, 0, len(bookings))
	for _, b := range bookings {
		respBookings = append(respBookings, &Booking{
			WorkshopID:     b.WorkshopID,
			ClientID:       b.ClientID,
			BeginAt:        b.BeginAt.Format("02-01-2006 15:04"),
			EndAt:          b.EndAt.Format("02-01-2006 15:04"),
			ClientTimezone: b.ClientTimezone.String(),
		})
	}

	c.JSON(200, ListBookingsResponse{Bookings: respBookings})

}
