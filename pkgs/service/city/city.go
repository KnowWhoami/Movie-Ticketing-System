package city

import (
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/pkgs/models"
)

type Service struct {
	db     *gorm.DB
	logger log.Logger
}

func NewService(db *gorm.DB, logger log.Logger) *Service {
	return &Service{db: db, logger: logger}
}

func (s *Service) AddCity(inp *AddCityInput) (*AddCityOutput, error) {
	c := models.City{
		Name:    inp.Name,
		ZipCode: inp.ZipCode,
	}
	if err := s.db.Create(&c).Error; err != nil {
		_ = s.logger.Log("method", "AddCity", "err", err)
		return nil, err
	}
	_ = s.logger.Log("method", "AddCity", "city_id", c.ID)
	return &AddCityOutput{City: c}, nil
}

func (s *Service) ListCities() (*ListCitiesOutput, error) {
	var cities []models.City
	if err := s.db.Find(&cities).Error; err != nil {
		_ = s.logger.Log("method", "ListCities", "err", err)
		return nil, err
	}
	return &ListCitiesOutput{Cities: cities}, nil
}
