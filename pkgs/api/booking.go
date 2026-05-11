package api

import (
	"net/http"
	"strconv"

	jsonHelper "KnowWhoami/movie-ticketing/internal/json"
	"KnowWhoami/movie-ticketing/pkgs/models"
	"KnowWhoami/movie-ticketing/pkgs/service/booking"
)

func (h *Handler) BookSeats(request *http.Request, writer http.ResponseWriter) {
	var inp booking.BookSeatsInput
	if err := ValidateContract(&inp, request, writer, h.db); err != nil {
		return
	}
	inp.UserID = TokenPayloadFromContext(request.Context()).UserID

	out, err := h.svc.booking.BookSeats(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "BookSeats", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	_ = h.logger.Log("handler", "BookSeats", "booking_id", out.Booking.ID)
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusCreated)
}

func (h *Handler) ListBookings(request *http.Request, writer http.ResponseWriter) {
	payload := TokenPayloadFromContext(request.Context())
	q := request.URL.Query()

	inp := booking.ListBookingsInput{
		CallerID:   payload.UserID,
		CallerType: models.UserType(payload.UserType),
	}
	if v := q.Get("movie_show_id"); v != "" {
		inp.MovieShowID, _ = strconv.Atoi(v)
	}
	if v := q.Get("user_id"); v != "" {
		inp.UserID, _ = strconv.Atoi(v)
	}

	out, err := h.svc.booking.ListBookings(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "ListBookings", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) ListMyBookings(request *http.Request, writer http.ResponseWriter) {
	userID := TokenPayloadFromContext(request.Context()).UserID

	out, err := h.svc.booking.ListMyBookings(userID)
	if err != nil {
		_ = h.logger.Log("handler", "ListMyBookings", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}
