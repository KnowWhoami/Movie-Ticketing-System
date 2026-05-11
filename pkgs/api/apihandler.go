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
	"KnowWhoami/movie-ticketing/pkgs/service/city"
	"KnowWhoami/movie-ticketing/pkgs/service/movie"
	userSvc "KnowWhoami/movie-ticketing/pkgs/service/user"
)

type Handler struct {
	db     *gorm.DB
	logger log.Logger

	// services
	svc *HandlerServices
}

type HandlerServices struct {
	city    *city.Service
	cinema  *cinema.Service
	movie   *movie.Service
	booking *booking.Service
	user    *userSvc.Service
}

func NewAPIHandler(db *gorm.DB, logger log.Logger) *Handler {
	return &Handler{
		db:     db,
		logger: logger,
		svc: &HandlerServices{
			city: city.NewService(db,
				log.With(logger, "service", "city")),
			cinema: cinema.NewService(db,
				cache.NewCache(cache.InMemoryCache), log.With(logger, "service", "cinema")),
			movie: movie.NewService(db,
				cache.NewCache(cache.InMemoryCache), log.With(logger, "service", "movie")),
			booking: booking.NewService(db,
				cache.NewCache(cache.InMemoryCache), log.With(logger, "service", "booking")),
			user: userSvc.NewService(db,
				log.With(logger, "service", "user")),
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
