package booking

import (
	"time"

	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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

func (s *Service) ListBookings() (*ListBookingsOutput, error) {
	if cachedValue, ok := s.cache.Get("ListBookingsOutput"); ok {
		return cachedValue.(*ListBookingsOutput), nil
	}
	var bookings []models.Booking

	// Get all records
	result := s.db.Preload(clause.Associations).Find(&bookings)
	if result.Error != nil {
		_ = s.logger.Log("method", "ListBookings", "err", result.Error)
		return nil, result.Error
	}
	out := &ListBookingsOutput{Bookings: bookings}
	s.cache.Set("ListBookingsOutput", out, time.Duration(10*time.Minute))
	return out, nil
}

func (s *Service) BookSeats(inp *BookSeatsInput) (*BookSeatsOutput, error) {
	var booking models.Booking
	var bookingErr error

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		booking = models.Booking{
			SeatCount:   len(inp.SeatNumbers),
			Status:      models.BookingPending,
			UserID:      inp.UserID,
			MovieShowID: inp.ShowID,
		}
		if result := tx.Create(&booking); result.Error != nil {
			return result.Error
		}

		if result := tx.Preload("MovieShow").Find(&booking, booking.ID); result.Error != nil {
			return result.Error
		}

		bookingErr = booking.BookSeats(tx, inp.SeatNumbers, inp.SeatType)
		if bookingErr != nil {
			if err := booking.Fail(); err != nil {
				return err
			}
			_ = s.logger.Log("method", "BookSeats", "booking_id", booking.ID, "err", bookingErr)
		} else {
			if err := booking.Confirm(); err != nil {
				return err
			}
		}

		if result := tx.Save(&booking); result.Error != nil {
			return result.Error
		}
		// Return nil so the transaction commits in both cases, preserving the
		// CONFIRMED or FAILED record for audit.
		return nil
	})

	s.cache.Delete("ListBookingsOutput")

	if txErr != nil {
		_ = s.logger.Log("method", "BookSeats", "err", txErr)
		return nil, txErr
	}

	if bookingErr != nil {
		return &BookSeatsOutput{Booking: booking}, bookingErr
	}

	if result := s.db.Preload("Seats").
		Preload("MovieShow").
		Preload("User").
		Find(&booking, booking.ID); result.Error != nil {
		return &BookSeatsOutput{Booking: booking}, result.Error
	}

	_ = s.logger.Log("method", "BookSeats", "booking_id", booking.ID, "status", booking.Status)
	return &BookSeatsOutput{Booking: booking}, nil
}
