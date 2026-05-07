package cinema

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

func (s *Service) AddCinema(inp *AddCinemaInput) (*AddCinemaOutput, error) {
	c := models.Cinema{
		Name:   inp.CinemaName,
		CityID: inp.CityID,
	}
	result := s.db.Create(&c)
	if result.Error != nil {
		_ = s.logger.Log("method", "AddCinema", "err", result.Error)
		return nil, result.Error
	}
	// added successfully
	s.cache.Delete("ListCinemasOutput")
	_ = s.logger.Log("method", "AddCinema", "cinema_id", c.ID)
	return &AddCinemaOutput{Cinema: c}, nil
}

func (s *Service) ListCinemas() (*ListCinemasOutput, error) {
	if cachedValue, ok := s.cache.Get("ListCinemasOutput"); ok {
		return cachedValue.(*ListCinemasOutput), nil
	}
	var cinemas []models.Cinema

	// Get all records
	result := s.db.Preload("CinemaScreens.CinemaSeats").Preload(clause.Associations).Find(&cinemas)
	if result.Error != nil {
		_ = s.logger.Log("method", "ListCinemas", "err", result.Error)
		return nil, result.Error
	}
	out := &ListCinemasOutput{Cinemas: cinemas}
	s.cache.Set("ListCinemasOutput", out, time.Duration(10*time.Minute))
	return out, nil
}

func (s *Service) AddCinemaScreen(inp *AddCinemaScreenInput) (*AddCinemaScreenOutput, error) {
	var out AddCinemaScreenOutput

	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		screen := models.CinemaScreen{
			Name:     inp.ScreenName,
			CinemaID: inp.CinemaID,
		}
		if result := tx.Create(&screen); result.Error != nil {
			return result.Error
		}

		if len(inp.Seats) > 0 {
			seats := make([]models.CinemaSeat, len(inp.Seats))
			for i, seat := range inp.Seats {
				seats[i] = models.CinemaSeat{
					SeatNumber:     seat.SeatNumber,
					Type:           seat.SeatType,
					CinemaScreenID: screen.ID,
				}
			}
			if result := tx.Create(&seats); result.Error != nil {
				return result.Error
			}
		}

		if result := tx.Preload("CinemaSeats").Preload("Cinema").Find(&screen, screen.ID); result.Error != nil {
			return result.Error
		}
		out.CinemaScreen = screen
		return nil
	})

	if txErr != nil {
		_ = s.logger.Log("method", "AddCinemaScreen", "err", txErr)
		return nil, txErr
	}

	s.cache.Delete("ListCinemasOutput")
	_ = s.logger.Log("method", "AddCinemaScreen", "screen_id", out.CinemaScreen.ID)
	return &out, nil
}
