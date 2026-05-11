package api

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	jsonHelper "KnowWhoami/movie-ticketing/internal/json"
	"KnowWhoami/movie-ticketing/pkgs/models"
	"KnowWhoami/movie-ticketing/pkgs/service/movie"
)

func (h *Handler) AddMovie(request *http.Request, writer http.ResponseWriter) {
	var inp movie.AddMovieInput
	if err := ValidateContract(&inp, request, writer, h.db); err != nil {
		return
	}

	out, err := h.svc.movie.AddMovie(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "AddMovie", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	_ = h.logger.Log("handler", "AddMovie", "movie_id", out.Movie.ID)
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusCreated)
}

func (h *Handler) ListMovies(request *http.Request, writer http.ResponseWriter) {
	inp := movie.ListMoviesInput{
		Name: request.URL.Query().Get("name"),
	}
	out, err := h.svc.movie.ListMovies(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "ListMovies", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) GetMovieByID(request *http.Request, writer http.ResponseWriter) {
	params := httprouter.ParamsFromContext(request.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id == 0 {
		resp := ErrorResponse("invalid movie id", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}

	out, err := h.svc.movie.GetMovieByID(id)
	if err != nil {
		_ = h.logger.Log("handler", "GetMovieByID", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusNotFound)
		jsonHelper.WriteResult(&resp, writer, http.StatusNotFound)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) UpdateMovie(request *http.Request, writer http.ResponseWriter) {
	params := httprouter.ParamsFromContext(request.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id == 0 {
		resp := ErrorResponse("invalid movie id", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}

	var inp movie.UpdateMovieInput
	if err := ValidateContract(&inp, request, writer, h.db); err != nil {
		return
	}
	inp.ID = id

	out, err := h.svc.movie.UpdateMovie(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "UpdateMovie", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) AddShow(request *http.Request, writer http.ResponseWriter) {
	var inp movie.AddMovieShowInput
	if err := ValidateContract(&inp, request, writer, h.db); err != nil {
		return
	}
	inp.CinemaOwnerID = TokenPayloadFromContext(request.Context()).UserID

	out, err := h.svc.movie.AddMovieShow(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "AddShow", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	_ = h.logger.Log("handler", "AddShow", "show_id", out.Show.ID)
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusCreated)
}

func (h *Handler) ListShows(request *http.Request, writer http.ResponseWriter) {
	q := request.URL.Query()
	inp := movie.ListShowsInput{}
	if v := q.Get("movie_id"); v != "" {
		inp.MovieID, _ = strconv.Atoi(v)
	}
	if v := q.Get("cinema_screen_id"); v != "" {
		inp.CinemaScreenID, _ = strconv.Atoi(v)
	}
	if v := q.Get("cinema_id"); v != "" {
		inp.CinemaID, _ = strconv.Atoi(v)
	}

	out, err := h.svc.movie.ListShows(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "ListShows", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) GetShowByID(request *http.Request, writer http.ResponseWriter) {
	params := httprouter.ParamsFromContext(request.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id == 0 {
		resp := ErrorResponse("invalid show id", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}

	out, err := h.svc.movie.GetShowByID(id)
	if err != nil {
		_ = h.logger.Log("handler", "GetShowByID", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusNotFound)
		jsonHelper.WriteResult(&resp, writer, http.StatusNotFound)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) ListShowSeats(request *http.Request, writer http.ResponseWriter) {
	params := httprouter.ParamsFromContext(request.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id == 0 {
		resp := ErrorResponse("invalid show id", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}

	q := request.URL.Query()
	inp := movie.ListShowSeatsInput{ShowID: id}

	if v := q.Get("type"); v != "" {
		inp.Type = models.SeatType(v)
	}
	if v := q.Get("available"); v != "" {
		b := v == "true"
		inp.Available = &b
	}

	out, err := h.svc.movie.ListShowSeats(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "ListShowSeats", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) CancelShow(request *http.Request, writer http.ResponseWriter) {
	params := httprouter.ParamsFromContext(request.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id == 0 {
		resp := ErrorResponse("invalid show id", http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusBadRequest)
		return
	}

	inp := movie.CancelShowInput{
		ShowID:        id,
		CinemaOwnerID: TokenPayloadFromContext(request.Context()).UserID,
	}

	out, err := h.svc.movie.CancelShow(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "CancelShow", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}
