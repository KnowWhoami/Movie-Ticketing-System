package city

import (
	"fmt"

	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/pkgs/models"
)

type AddCityInput struct {
	Name    string `json:"name"`
	ZipCode string `json:"zip_code"`
}

func (a *AddCityInput) Validate(db *gorm.DB) error {
	if a.Name == "" || a.ZipCode == "" {
		return fmt.Errorf("name and zip_code are required")
	}
	var existing models.City
	if err := db.Where("name = ?", a.Name).First(&existing).Error; err == nil {
		return fmt.Errorf("city %q already exists", a.Name)
	}
	return nil
}

type AddCityOutput struct {
	City models.City `json:"city"`
}

type ListCitiesOutput struct {
	Cities []models.City `json:"cities"`
}
