package movie

import (
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

func (s *Service) AddMovie(inp *AddMovieInput) (*AddMovieOutput, error) {
	movie := models.Movie{
		Name:        inp.Name,
		Description: inp.Description,
		Duration:    inp.Duration,
	}

	if result := s.db.Create(&movie); result.Error != nil {
		_ = s.logger.Log("method", "AddMovie", "err", result.Error)
		return nil, result.Error
	}
	_ = s.logger.Log("method", "AddMovie", "movie_id", movie.ID)
	return &AddMovieOutput{Movie: movie}, nil
}

func (s *Service) AddMovieShow(inp *AddMovieShowInput) (*AddMovieShowOutput, error) {
	show := models.MovieShow{
		StartTime:      inp.StartTime,
		EndTime:        inp.EndTime,
		MovieID:        inp.MovieID,
		CinemaScreenID: inp.CinemaScreenID,
	}

	// Explicit transaction is required so that the SELECT FOR UPDATE in
	// CheckOverlap and the INSERT share the same transaction boundary. Without
	// it the lock would release on autocommit before the INSERT, leaving the
	// race open.
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&show).Error
	}); err != nil {
		_ = s.logger.Log("method", "AddMovieShow", "err", err)
		return nil, err
	}

	out := &AddMovieShowOutput{}
	if result := s.db.Preload("Movie").Preload("Bookings").Preload("Seats").First(&show, show.ID); result.Error != nil {
		_ = s.logger.Log("method", "AddMovieShow", "err", result.Error)
		return nil, result.Error
	}
	out.Show = show
	_ = s.logger.Log("method", "AddMovieShow", "show_id", show.ID)
	return out, nil
}

func (s *Service) GetMovieShow(inp *GetMovieShowInput) (*GetMovieShowOutput, error) {
	var show models.MovieShow
	if result := s.db.Preload(clause.Associations).Find(&show, inp.ShowID); result.Error != nil {
		_ = s.logger.Log("method", "GetMovieShow", "err", result.Error)
		return nil, result.Error
	}
	return &GetMovieShowOutput{Show: show}, nil
}
