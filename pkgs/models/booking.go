package models

import (
	"fmt"

	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingConfirmed BookingStatus = "CONFIRMED"
	BookingCancelled BookingStatus = "CANCELLED"
	BookingFailed    BookingStatus = "FAILED"
	BookingPending   BookingStatus = "PENDING"
)

type Booking struct {
	Model
	Status BookingStatus `json:"status"`

	// relations
	UserID      int        `json:"-"`
	User        User       `json:"user"`
	MovieShowID int        `json:"-"`
	MovieShow   MovieShow  `json:"movie_show"`
	Seats       []MovieShowSeat `json:"seats" gorm:"foreignKey:BookingID"`
}

func (b *Booking) BeforeSave(db *gorm.DB) (err error) {
	b.MovieShow = MovieShow{}
	return
}

// BookSeats claims the given show_seat IDs for this booking.
// Uses a conditional UPDATE so concurrent requests race on the DB constraint.
func (b *Booking) BookSeats(db *gorm.DB, showSeatIDs []int) error {
	result := db.Model(&MovieShowSeat{}).
		Where("id IN ? AND booking_id IS NULL", showSeatIDs).
		Update("booking_id", b.ID)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != int64(len(showSeatIDs)) {
		// revert any seats that were successfully claimed in this attempt
		revert := db.Model(&MovieShowSeat{}).
			Where("id IN ? AND booking_id = ?", showSeatIDs, b.ID).
			Update("booking_id", nil)
		errMsg := "some seats are already booked"
		if revert.Error != nil {
			errMsg += ": revert error: " + revert.Error.Error()
		}
		return fmt.Errorf("%s", errMsg)
	}
	return nil
}

func (b *Booking) Fail() error {
	if b.Status != BookingPending {
		return fmt.Errorf("cannot fail booking with status %s", b.Status)
	}
	b.Status = BookingFailed
	return nil
}

func (b *Booking) Confirm() error {
	if b.Status != BookingPending {
		return fmt.Errorf("cannot confirm booking with status %s", b.Status)
	}
	b.Status = BookingConfirmed
	return nil
}
