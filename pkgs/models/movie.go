package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Movie struct {
	Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Duration    int    `json:"duration"` // minutes

	// relations
	MovieShows []MovieShow `json:"shows" gorm:"foreignKey:MovieID"`
}

type MovieShow struct {
	Model
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	IsCancelled bool      `json:"is_cancelled" gorm:"default:false"`

	// relations
	CinemaScreenID int `json:"-"`
	MovieID        int `json:"-"`
	CinemaScreen   CinemaScreen
	Movie          Movie
	Bookings       []Booking `json:"bookings" gorm:"foreignKey:MovieShowID"`
}

func (ms *MovieShow) BeforeCreate(db *gorm.DB) (err error) {
	return ms.CheckOverlap(db)
}

func (ms *MovieShow) CheckOverlap(db *gorm.DB) error {
	var shows []*MovieShow

	// FOR UPDATE locks the matched rows (and the gap when there are none) so
	// that a concurrent transaction on the same screen must wait until this
	// transaction commits before it can read. Combined with the explicit
	// transaction in AddMovieShow, this makes the check-then-insert atomic.
	if result := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("cinema_screen_id = ? AND start_time < ? AND end_time > ?",
			ms.CinemaScreenID, ms.EndTime, ms.StartTime).
		Find(&shows); result.Error != nil {
		return result.Error
	}
	if len(shows) > 0 {
		return fmt.Errorf("the show has an overlap with %d shows", len(shows))
	}
	return nil
}

func (ms *MovieShow) AfterCreate(db *gorm.DB) error {
	return ms.GenerateShowSeats(db)
}

func (ms *MovieShow) GenerateShowSeats(db *gorm.DB) error {
	var seats []*CinemaSeat
	if err := db.Where("cinema_screen_id = ?", ms.CinemaScreenID).Find(&seats).Error; err != nil {
		return err
	}

	showSeats := make([]*MovieShowSeat, len(seats))
	for i, seat := range seats {
		showSeats[i] = &MovieShowSeat{
			MovieShowID:  ms.ID,
			CinemaSeatID: seat.ID,
		}
	}
	return db.Create(showSeats).Error
}

type MovieShowSeat struct {
	Model
	MovieShowID  int  `json:"-" gorm:"uniqueIndex:unique_seat_per_show"`
	CinemaSeatID int  `json:"-" gorm:"uniqueIndex:unique_seat_per_show"`
	BookingID    *int `json:"-"`

	// relations
	MovieShow  MovieShow  `json:"-"`
	CinemaSeat CinemaSeat `json:"cinema_seat"`
	Booking    *Booking   `json:"-"`
}
