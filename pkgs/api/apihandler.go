package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/internal/cache"
	jsonHelper "KnowWhoami/movie-ticketing/internal/json"
	"KnowWhoami/movie-ticketing/pkgs/service/booking"
	"KnowWhoami/movie-ticketing/pkgs/service/cinema"
	"KnowWhoami/movie-ticketing/pkgs/service/movie"
)

type Handler struct {
	db     *gorm.DB
	logger log.Logger

	// services
	svc *HandlerServices
}

type HandlerServices struct {
	cinema  *cinema.Service
	movie   *movie.Service
	booking *booking.Service
}

func NewAPIHandler(db *gorm.DB, logger log.Logger) *Handler {
	return &Handler{
		db:     db,
		logger: logger,
		svc: &HandlerServices{
			cinema: cinema.NewService(db,
				cache.NewCache(cache.InMemoryCache), log.With(logger, "service", "cinema")),
			movie: movie.NewService(db,
				cache.NewCache(cache.InMemoryCache), log.With(logger, "service", "movie")),
			booking: booking.NewService(db,
				cache.NewCache(cache.InMemoryCache), log.With(logger, "service", "booking")),
		},
	}
}

type HandlerFunc func(request *http.Request, writer http.ResponseWriter)

type Contract interface {
	Validate(db *gorm.DB) error
}

func ValidateContract(c Contract, request *http.Request, writer http.ResponseWriter, db *gorm.DB) error {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(c)
	if err == nil {
		err = c.Validate(db)
	}
	if err != nil {
		resp := ErrorResponse(err.Error(), http.StatusBadRequest)
		jsonHelper.WriteResult(&resp, writer, http.StatusUnprocessableEntity)
		return err
	}
	return nil
}
