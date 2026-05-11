package booking

import (
	"fmt"

	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/pkgs/models"
)

type BookSeatsInput struct {
	MovieShowID int   `json:"movie_show_id"`
	ShowSeatIDs []int `json:"show_seat_ids"`
	UserID      int   `json:"-"` // set from JWT by handler
}

func (b *BookSeatsInput) Validate(db *gorm.DB) error {
	if b.MovieShowID == 0 {
		return fmt.Errorf("movie_show_id is required")
	}
	if len(b.ShowSeatIDs) == 0 {
		return fmt.Errorf("show_seat_ids must not be empty")
	}
	var show models.MovieShow
	if err := db.First(&show, b.MovieShowID).Error; err != nil {
		return fmt.Errorf("show not found")
	}
	if show.IsCancelled {
		return fmt.Errorf("show has been cancelled")
	}
	// verify every seat ID belongs to this show
	var count int64
	db.Model(&models.MovieShowSeat{}).
		Where("id IN ? AND movie_show_id = ?", b.ShowSeatIDs, b.MovieShowID).
		Count(&count)
	if int(count) != len(b.ShowSeatIDs) {
		return fmt.Errorf("one or more seat IDs do not belong to this show")
	}
	return nil
}

type BookSeatsOutput struct {
	Booking models.Booking `json:"booking"`
}

type ListBookingsInput struct {
	MovieShowID int
	UserID      int
	CallerID    int
	CallerType  models.UserType
}

type ListBookingsOutput struct {
	Bookings []models.Booking `json:"bookings"`
}

type ListMyBookingsOutput struct {
	Bookings []models.Booking `json:"bookings"`
}
