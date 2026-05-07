package api

import (
	"net/http"

	jsonHelper "KnowWhoami/movie-ticketing/internal/json"
	"KnowWhoami/movie-ticketing/pkgs/service/movie"
)

func (h *Handler) GetShow(request *http.Request, writer http.ResponseWriter) {
	var inp movie.GetMovieShowInput
	if err := ValidateContract(&inp, request, writer, h.db); err != nil {
		return
	}

	out, err := h.svc.movie.GetMovieShow(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "GetShow", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	_ = h.logger.Log("handler", "GetShow", "show_id", inp.ShowID)
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

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
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}

func (h *Handler) AddShow(request *http.Request, writer http.ResponseWriter) {
	var inp movie.AddMovieShowInput
	if err := ValidateContract(&inp, request, writer, h.db); err != nil {
		return
	}

	out, err := h.svc.movie.AddMovieShow(&inp)
	if err != nil {
		_ = h.logger.Log("handler", "AddShow", "err", err)
		resp := ErrorResponse(err.Error(), http.StatusUnprocessableEntity)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return
	}
	_ = h.logger.Log("handler", "AddShow", "show_id", out.Show.ID)
	resp := SuccessResponse(out)
	jsonHelper.WriteResult(&resp, writer, http.StatusOK)
}
