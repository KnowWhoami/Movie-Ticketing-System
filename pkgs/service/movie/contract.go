package movie

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/pkgs/models"
)

type AddMovieInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Duration    int    `json:"duration"` // minutes
}

func (am *AddMovieInput) Validate(db *gorm.DB) error {
	if am.Name == "" || am.Description == "" || am.Duration == 0 {
		return fmt.Errorf("name, description and duration are mandatory fields")
	}
	return nil
}

type AddMovieOutput struct {
	Movie models.Movie `json:"movie"`
}

type ListMoviesInput struct {
	Name string
}

type ListMoviesOutput struct {
	Movies []models.Movie `json:"movies"`
}

type GetMovieByIDOutput struct {
	Movie models.Movie `json:"movie"`
}

type UpdateMovieInput struct {
	ID          int
	Name        string `json:"name"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
}

func (u *UpdateMovieInput) Validate(db *gorm.DB) error {
	return nil
}

type UpdateMovieOutput struct {
	Movie models.Movie `json:"movie"`
}

type AddMovieShowInput struct {
	MovieID        int       `json:"movie_id"`
	CinemaScreenID int       `json:"cinema_screen_id"`
	StartTime      time.Time `json:"start_time"`
	CinemaOwnerID  int       `json:"-"` // set from JWT by handler
}

func (ams *AddMovieShowInput) Validate(db *gorm.DB) error {
	if ams.MovieID == 0 || ams.CinemaScreenID == 0 || ams.StartTime.IsZero() {
		return fmt.Errorf("movie_id, cinema_screen_id and start_time are mandatory fields")
	}
	return nil
}

type AddMovieShowOutput struct {
	Show models.MovieShow `json:"show"`
}

type ListShowsInput struct {
	CinemaID       int
	CinemaScreenID int
	MovieID        int
}

type ListShowsOutput struct {
	Shows []models.MovieShow `json:"shows"`
}

type GetShowByIDOutput struct {
	Show models.MovieShow `json:"show"`
}

type CancelShowInput struct {
	ShowID        int
	CinemaOwnerID int // set from JWT by handler
}

type CancelShowOutput struct {
	Show models.MovieShow `json:"show"`
}

type ListShowSeatsInput struct {
	ShowID    int
	Type      models.SeatType
	Available *bool // nil = all, true = available only, false = booked only
}

// ShowSeatInfo is the response shape — derives "available" from BookingID instead of exposing the raw FK.
type ShowSeatInfo struct {
	ID         int    `json:"id"`
	SeatNumber int    `json:"seat_number"`
	Type       string `json:"type"`
	Available  bool   `json:"available"`
}

type ListShowSeatsOutput struct {
	Seats []ShowSeatInfo `json:"seats"`
}
