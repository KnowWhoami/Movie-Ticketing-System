package booking

import (
	"fmt"

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
	return &Service{db: db, cache: cache, logger: logger}
}

func (s *Service) BookSeats(inp *BookSeatsInput) (*BookSeatsOutput, error) {
	var booking models.Booking
	var bookingErr error

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		booking = models.Booking{
			Status:      models.BookingPending,
			UserID:      inp.UserID,
			MovieShowID: inp.MovieShowID,
		}
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		bookingErr = booking.BookSeats(tx, inp.ShowSeatIDs)
		if bookingErr != nil {
			_ = booking.Fail()
		} else {
			_ = booking.Confirm()
		}

		if err := tx.Save(&booking).Error; err != nil {
			return err
		}
		// commit in both cases — CONFIRMED and FAILED are both valid audit records
		return nil
	})

	if txErr != nil {
		_ = s.logger.Log("method", "BookSeats", "err", txErr)
		return nil, txErr
	}

	if bookingErr != nil {
		_ = s.logger.Log("method", "BookSeats", "booking_id", booking.ID, "err", bookingErr)
		return &BookSeatsOutput{Booking: booking}, bookingErr
	}

	if err := s.db.Preload("Seats.CinemaSeat").Preload("MovieShow").Preload("User").
		First(&booking, booking.ID).Error; err != nil {
		return &BookSeatsOutput{Booking: booking}, err
	}

	_ = s.logger.Log("method", "BookSeats", "booking_id", booking.ID, "status", booking.Status)
	return &BookSeatsOutput{Booking: booking}, nil
}

func (s *Service) ListBookings(inp *ListBookingsInput) (*ListBookingsOutput, error) {
	if inp.CallerType == models.UserTypeCinemaOwner {
		if inp.MovieShowID == 0 {
			return nil, fmt.Errorf("movie_show_id is required for cinema owners")
		}
		// verify the show belongs to the caller's cinema
		var show models.MovieShow
		if err := s.db.Preload("CinemaScreen.Cinema").First(&show, inp.MovieShowID).Error; err != nil {
			return nil, fmt.Errorf("show not found")
		}
		if show.CinemaScreen.Cinema.CinemaOwnerID != inp.CallerID {
			return nil, fmt.Errorf("show does not belong to your cinema")
		}
	}

	query := s.db.Preload("Seats.CinemaSeat").Preload("MovieShow").Preload("User")
	if inp.MovieShowID != 0 {
		query = query.Where("movie_show_id = ?", inp.MovieShowID)
	}
	if inp.UserID != 0 {
		query = query.Where("user_id = ?", inp.UserID)
	}

	var bookings []models.Booking
	if err := query.Find(&bookings).Error; err != nil {
		_ = s.logger.Log("method", "ListBookings", "err", err)
		return nil, err
	}
	return &ListBookingsOutput{Bookings: bookings}, nil
}

func (s *Service) ListMyBookings(userID int) (*ListMyBookingsOutput, error) {
	var bookings []models.Booking
	err := s.db.Preload("Seats.CinemaSeat").Preload("MovieShow").
		Where("user_id = ?", userID).Find(&bookings).Error
	if err != nil {
		_ = s.logger.Log("method", "ListMyBookings", "err", err)
		return nil, err
	}
	return &ListMyBookingsOutput{Bookings: bookings}, nil
}
