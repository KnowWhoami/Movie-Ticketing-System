package api

import (
	"net/http"

	jsonHelper "KnowWhoami/movie-ticketing/internal/json"
	"KnowWhoami/movie-ticketing/pkgs/service/cinema"
)

func (h *Handler) ListCinemas(request *http.Request, writer http.ResponseWriter) {
	out, err := h.svc.cinema.ListCinemas()
	if err != nil {
		_ = h.logger.Log("handler", "ListCinemas", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) AddCinema(request *http.Request, writer http.ResponseWriter) {
	var inp cinema.AddCinemaInput
	if err := ValidateContract(&inp, request, writer, h.db); err != nil {
		return
	}

	out, err := h.svc.cinema.AddCinema(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "AddCinema", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	_ = h.logger.Log("handler", "AddCinema", "cinema_id", out.Cinema.ID)
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) AddCinemaScreen(request *http.Request, writer http.ResponseWriter) {
	var inp cinema.AddCinemaScreenInput
	if err := ValidateContract(&inp, request, writer, h.db); err != nil {
		return
	}

	out, err := h.svc.cinema.AddCinemaScreen(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "AddCinemaScreen", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	_ = h.logger.Log("handler", "AddCinemaScreen", "screen_id", out.CinemaScreen.ID)
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}
