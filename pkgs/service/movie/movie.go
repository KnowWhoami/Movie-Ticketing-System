package movie

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/internal/cache"
	"KnowWhoami/movie-ticketing/pkgs/models"
)

type Service struct {
	db     *gorm.DB
	cache  cache.Cache
	logger log.Logger
}

func NewService(db *gorm.DB, cache cache.Cache, logger log.Logger) *Service {
	return &Service{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

func (s *Service) AddMovie(inp *AddMovieInput) (*AddMovieOutput, error) {
	movie := models.Movie{
		Name:        inp.Name,
		Description: inp.Description,
		Duration:    inp.Duration,
	}

	if err := s.db.Create(&movie).Error; err != nil {
		_ = s.logger.Log("method", "AddMovie", "err", err)
		return nil, err
	}
	_ = s.logger.Log("method", "AddMovie", "movie_id", movie.ID)
	return &AddMovieOutput{Movie: movie}, nil
}

func (s *Service) ListMovies(inp *ListMoviesInput) (*ListMoviesOutput, error) {
	query := s.db.Model(&models.Movie{})
	if inp.Name != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+inp.Name+"%")
	}
	var movies []models.Movie
	if err := query.Find(&movies).Error; err != nil {
		_ = s.logger.Log("method", "ListMovies", "err", err)
		return nil, err
	}
	return &ListMoviesOutput{Movies: movies}, nil
}

func (s *Service) GetMovieByID(id int) (*GetMovieByIDOutput, error) {
	var movie models.Movie
	if err := s.db.Preload("MovieShows").First(&movie, id).Error; err != nil {
		return nil, fmt.Errorf("movie not found")
	}
	return &GetMovieByIDOutput{Movie: movie}, nil
}

func (s *Service) UpdateMovie(inp *UpdateMovieInput) (*UpdateMovieOutput, error) {
	var movie models.Movie
	if err := s.db.First(&movie, inp.ID).Error; err != nil {
		return nil, fmt.Errorf("movie not found")
	}
	if inp.Name != "" {
		movie.Name = inp.Name
	}
	if inp.Description != "" {
		movie.Description = inp.Description
	}
	if inp.Duration != 0 {
		movie.Duration = inp.Duration
	}
	if err := s.db.Save(&movie).Error; err != nil {
		_ = s.logger.Log("method", "UpdateMovie", "err", err)
		return nil, err
	}
	_ = s.logger.Log("method", "UpdateMovie", "movie_id", movie.ID)
	return &UpdateMovieOutput{Movie: movie}, nil
}

func (s *Service) AddMovieShow(inp *AddMovieShowInput) (*AddMovieShowOutput, error) {
	// verify screen belongs to the requesting cinema owner
	var screen models.CinemaScreen
	if err := s.db.Preload("Cinema").First(&screen, inp.CinemaScreenID).Error; err != nil {
		return nil, fmt.Errorf("cinema screen not found")
	}
	if screen.Cinema.CinemaOwnerID != inp.CinemaOwnerID {
		return nil, fmt.Errorf("cinema screen does not belong to your cinema")
	}

	var movie models.Movie
	if err := s.db.First(&movie, inp.MovieID).Error; err != nil {
		return nil, fmt.Errorf("movie not found")
	}

	show := models.MovieShow{
		StartTime:      inp.StartTime,
		EndTime:        inp.StartTime.Add(time.Duration(movie.Duration) * time.Minute),
		MovieID:        inp.MovieID,
		CinemaScreenID: inp.CinemaScreenID,
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&show).Error
	}); err != nil {
		_ = s.logger.Log("method", "AddMovieShow", "err", err)
		return nil, err
	}

	if err := s.db.Preload("Movie").Preload("CinemaScreen").First(&show, show.ID).Error; err != nil {
		return nil, err
	}
	_ = s.logger.Log("method", "AddMovieShow", "show_id", show.ID)
	return &AddMovieShowOutput{Show: show}, nil
}

func (s *Service) ListShows(inp *ListShowsInput) (*ListShowsOutput, error) {
	query := s.db.Preload("Movie").Preload("CinemaScreen")

	if inp.MovieID != 0 {
		query = query.Where("movie_id = ?", inp.MovieID)
	}
	if inp.CinemaScreenID != 0 {
		query = query.Where("cinema_screen_id = ?", inp.CinemaScreenID)
	}
	if inp.CinemaID != 0 {
		query = query.Joins("JOIN cinema_screens ON cinema_screens.id = movie_shows.cinema_screen_id").
			Where("cinema_screens.cinema_id = ?", inp.CinemaID)
	}

	var shows []models.MovieShow
	if err := query.Find(&shows).Error; err != nil {
		_ = s.logger.Log("method", "ListShows", "err", err)
		return nil, err
	}
	return &ListShowsOutput{Shows: shows}, nil
}

func (s *Service) GetShowByID(id int) (*GetShowByIDOutput, error) {
	var show models.MovieShow
	if err := s.db.Preload("Movie").Preload("CinemaScreen").First(&show, id).Error; err != nil {
		return nil, fmt.Errorf("show not found")
	}
	return &GetShowByIDOutput{Show: show}, nil
}

func (s *Service) CancelShow(inp *CancelShowInput) (*CancelShowOutput, error) {
	var show models.MovieShow
	if err := s.db.Preload("CinemaScreen.Cinema").First(&show, inp.ShowID).Error; err != nil {
		return nil, fmt.Errorf("show not found")
	}
	if show.CinemaScreen.Cinema.CinemaOwnerID != inp.CinemaOwnerID {
		return nil, fmt.Errorf("you can only cancel shows for your own cinemas")
	}
	if show.IsCancelled {
		return nil, fmt.Errorf("show is already cancelled")
	}

	if err := s.db.Model(&show).Update("is_cancelled", true).Error; err != nil {
		_ = s.logger.Log("method", "CancelShow", "err", err)
		return nil, err
	}
	show.IsCancelled = true
	_ = s.logger.Log("method", "CancelShow", "show_id", show.ID)
	return &CancelShowOutput{Show: show}, nil
}

func (s *Service) ListShowSeats(inp *ListShowSeatsInput) (*ListShowSeatsOutput, error) {
	query := s.db.Model(&models.MovieShowSeat{}).
		Where("movie_show_seats.movie_show_id = ?", inp.ShowID).
		Preload("CinemaSeat")

	if inp.Type != "" {
		query = query.
			Joins("JOIN cinema_seats ON cinema_seats.id = movie_show_seats.cinema_seat_id").
			Where("cinema_seats.type = ?", inp.Type)
	}
	if inp.Available != nil {
		if *inp.Available {
			query = query.Where("movie_show_seats.booking_id IS NULL")
		} else {
			query = query.Where("movie_show_seats.booking_id IS NOT NULL")
		}
	}

	var seats []models.MovieShowSeat
	if err := query.Find(&seats).Error; err != nil {
		_ = s.logger.Log("method", "ListShowSeats", "err", err)
		return nil, err
	}

	infos := make([]ShowSeatInfo, len(seats))
	for i, seat := range seats {
		infos[i] = ShowSeatInfo{
			ID:         seat.ID,
			SeatNumber: seat.CinemaSeat.SeatNumber,
			Type:       string(seat.CinemaSeat.Type),
			Available:  seat.BookingID == nil,
		}
	}
	return &ListShowSeatsOutput{Seats: infos}, nil
}
