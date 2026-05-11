package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	jsonHelper "KnowWhoami/movie-ticketing/internal/json"
	"KnowWhoami/movie-ticketing/pkgs/models"
	"KnowWhoami/movie-ticketing/pkgs/service/cinema"
)

func (h *Handler) ListCinemas(request *http.Request, writer http.ResponseWriter) {
	q := request.URL.Query()
	inp := cinema.ListCinemasInput{
		Name: q.Get("name"),
	}
	if v := q.Get("city_id"); v != "" {
		inp.CityID, _ = strconv.Atoi(v)
	}

	out, err := h.svc.cinema.ListCinemas(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "ListCinemas", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) GetCinemasByOwnerID(request *http.Request, writer http.ResponseWriter) {
	params := httprouter.ParamsFromContext(request.Context())
	ownerID, err := strconv.Atoi(params.ByName("cinema_owner_id"))
	if err != nil || ownerID == 0 {
		resp := ErrorResponse("invalid cinema_owner_id", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}

	out, err := h.svc.cinema.GetMyCinemas(ownerID)
	if err != nil {
		_ = h.logger.Log("handler", "GetCinemasByOwnerID", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) GetCinemaByID(request *http.Request, writer http.ResponseWriter) {
	params := httprouter.ParamsFromContext(request.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id == 0 {
		resp := ErrorResponse("invalid cinema id", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}

	out, err := h.svc.cinema.GetCinemaByID(id)
	if err != nil {
		_ = h.logger.Log("handler", "GetCinemaByID", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusNotFound)
		jsonHelper.WriteResult(&resp, writer, http.StatusNotFound)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) GetMyCinemas(request *http.Request, writer http.ResponseWriter) {
	payload := TokenPayloadFromContext(request.Context())

	out, err := h.svc.cinema.GetMyCinemas(payload.UserID)
	if err != nil {
		_ = h.logger.Log("handler", "GetMyCinemas", "err", err)
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
	jsonHelper.WriteResult(&resp, writer, http.StatusCreated)
}

func (h *Handler) UpdateCinema(request *http.Request, writer http.ResponseWriter) {
	params := httprouter.ParamsFromContext(request.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id == 0 {
		resp := ErrorResponse("invalid cinema id", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}

	var inp cinema.UpdateCinemaInput
	if err := ValidateContract(&inp, request, writer, h.db); err != nil {
		return
	}
	inp.ID = id

	out, err := h.svc.cinema.UpdateCinema(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "UpdateCinema", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) ListSeats(request *http.Request, writer http.ResponseWriter) {
	params := httprouter.ParamsFromContext(request.Context())
	screenID, err := strconv.Atoi(params.ByName("screen_id"))
	if err != nil || screenID == 0 {
		resp := ErrorResponse("invalid screen_id", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}
	inp := cinema.ListSeatsInput{
		CinemaScreenID: screenID,
		Type:           models.SeatType(request.URL.Query().Get("type")),
	}
	out, err := h.svc.cinema.ListSeats(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "ListSeats", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) GetSeatByNumber(request *http.Request, writer http.ResponseWriter) {
	params := httprouter.ParamsFromContext(request.Context())
	screenID, err := strconv.Atoi(params.ByName("screen_id"))
	if err != nil || screenID == 0 {
		resp := ErrorResponse("invalid screen_id", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}
	seatNumber, err := strconv.Atoi(params.ByName("seat_number"))
	if err != nil || seatNumber == 0 {
		resp := ErrorResponse("invalid seat_number", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}
	out, err := h.svc.cinema.GetSeatByNumber(screenID, seatNumber)
	if err != nil {
		_ = h.logger.Log("handler", "GetSeatByNumber", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusNotFound)
		jsonHelper.WriteResult(&resp, writer, http.StatusNotFound)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) AddCinemaScreen(request *http.Request, writer http.ResponseWriter) {
	var inp cinema.AddCinemaScreenInput
	if err := ValidateContract(&inp, request, writer, h.db); err != nil {
		return
	}
	inp.CinemaOwnerID = TokenPayloadFromContext(request.Context()).UserID

	out, err := h.svc.cinema.AddCinemaScreen(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "AddCinemaScreen", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	_ = h.logger.Log("handler", "AddCinemaScreen", "screen_id", out.CinemaScreen.ID)
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusCreated)
}

func (h *Handler) ListScreens(request *http.Request, writer http.ResponseWriter) {
	q := request.URL.Query()
	inp := cinema.ListScreensInput{
		Name: q.Get("name"),
	}
	if v := q.Get("cinema_id"); v != "" {
		inp.CinemaID, _ = strconv.Atoi(v)
	}
	out, err := h.svc.cinema.ListScreens(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "ListScreens", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}
