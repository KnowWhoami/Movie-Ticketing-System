package api

import (
	"net/http"

	jsonHelper "KnowWhoami/movie-ticketing/internal/json"
	"KnowWhoami/movie-ticketing/pkgs/service/city"
)

func (h *Handler) AddCity(request *http.Request, writer http.ResponseWriter) {
	var inp city.AddCityInput
	if err := ValidateContract(&inp, request, writer, h.db); err != nil {
		return
	}
	out, err := h.svc.city.AddCity(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "AddCity", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusCreated)
}

func (h *Handler) ListCities(request *http.Request, writer http.ResponseWriter) {
	out, err := h.svc.city.ListCities()
	if err != nil {
		_ = h.logger.Log("handler", "ListCities", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}
