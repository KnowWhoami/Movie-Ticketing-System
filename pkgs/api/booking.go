package api

import (
	"net/http"

	jsonHelper "KnowWhoami/movie-ticketing/internal/json"
	"KnowWhoami/movie-ticketing/pkgs/service/booking"
)

func (h *Handler) ListBookings(request *http.Request, writer http.ResponseWriter) {
	out, err := h.svc.booking.ListBookings()
	if err != nil {
		_ = h.logger.Log("handler", "ListBookings", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) BookSeats(request *http.Request, writer http.ResponseWriter) {
	var inp booking.BookSeatsInput
	if err := ValidateContract(&inp, request, writer, h.db); err != nil {
		return
	}

	out, err := h.svc.booking.BookSeats(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "BookSeats", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	_ = h.logger.Log("handler", "BookSeats", "booking_id", out.Booking.ID)
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}
