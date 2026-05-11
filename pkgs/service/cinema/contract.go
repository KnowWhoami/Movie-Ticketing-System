package cinema

import (
	"fmt"

	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/pkgs/models"
)

type AddCinemaInput struct {
	CinemaName    string `json:"cinema_name"`
	CityID        int    `json:"city_id"`
	CinemaOwnerID int    `json:"cinema_owner_id"`
}

func (ac *AddCinemaInput) Validate(db *gorm.DB) error {
	if ac.CinemaName == "" || ac.CityID == 0 || ac.CinemaOwnerID == 0 {
		return fmt.Errorf("cinema_name, city_id, and cinema_owner_id are required")
	}
	var city models.City
	if err := db.First(&city, ac.CityID).Error; err != nil {
		return fmt.Errorf("city not found")
	}
	var owner models.User
	if err := db.First(&owner, ac.CinemaOwnerID).Error; err != nil {
		return fmt.Errorf("cinema owner not found")
	}
	if owner.UserType != models.UserTypeCinemaOwner {
		return fmt.Errorf("user %d is not a cinema owner", ac.CinemaOwnerID)
	}
	return nil
}

type AddCinemaOutput struct {
	Cinema models.Cinema `json:"cinema"`
}

type ListCinemasInput struct {
	Name   string
	CityID int
}

type ListCinemasOutput struct {
	Cinemas []models.Cinema `json:"cinemas"`
}

type GetCinemaByIDOutput struct {
	Cinema models.Cinema `json:"cinema"`
}

type UpdateCinemaInput struct {
	ID            int
	CinemaName    string `json:"cinema_name"`
	CityID        int    `json:"city_id"`
	CinemaOwnerID int    `json:"cinema_owner_id"`
}

func (u *UpdateCinemaInput) Validate(db *gorm.DB) error {
	if u.CityID != 0 {
		var city models.City
		if err := db.First(&city, u.CityID).Error; err != nil {
			return fmt.Errorf("city not found")
		}
	}
	if u.CinemaOwnerID != 0 {
		var owner models.User
		if err := db.First(&owner, u.CinemaOwnerID).Error; err != nil {
			return fmt.Errorf("cinema owner not found")
		}
		if owner.UserType != models.UserTypeCinemaOwner {
			return fmt.Errorf("user %d is not a cinema owner", u.CinemaOwnerID)
		}
	}
	return nil
}

type UpdateCinemaOutput struct {
	Cinema models.Cinema `json:"cinema"`
}

type AddCinemaScreenInput struct {
	CinemaID      int                     `json:"cinema_id"`
	ScreenName    string                  `json:"screen_name"`
	Seats         map[models.SeatType]int `json:"seats"` // e.g. {"PREMIUM": 30, "RECLINER": 10}
	CinemaOwnerID int                     `json:"-"`     // set from JWT by handler
}

func (acs *AddCinemaScreenInput) Validate(db *gorm.DB) error {
	if acs.CinemaID == 0 || acs.ScreenName == "" {
		return fmt.Errorf("cinema_id and screen_name are required")
	}
	if len(acs.Seats) == 0 {
		return fmt.Errorf("at least one seat type with a count must be provided")
	}
	validTypes := map[models.SeatType]bool{
		models.Recliner: true,
		models.Premium:  true,
		models.FrontRow: true,
		models.Balcony:  true,
	}
	for seatType, count := range acs.Seats {
		if !validTypes[seatType] {
			return fmt.Errorf("invalid seat type %q", seatType)
		}
		if count <= 0 {
			return fmt.Errorf("seat count for %q must be greater than 0", seatType)
		}
	}
	var cinema models.Cinema
	if err := db.First(&cinema, acs.CinemaID).Error; err != nil {
		return fmt.Errorf("cinema not found")
	}
	return nil
}

// CinemaScreenSummary is the response shape for screen endpoints — counts per type, not individual seats.
type CinemaScreenSummary struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	CinemaID     int            `json:"cinema_id"`
	SeatsSummary map[string]int `json:"seats_summary"`
}

type AddCinemaScreenOutput struct {
	CinemaScreen CinemaScreenSummary `json:"cinema_screen"`
}

type ListScreensInput struct {
	CinemaID int
	Name     string
}

type ListScreensOutput struct {
	CinemaScreens []CinemaScreenSummary `json:"cinema_screens"`
}

type ListSeatsInput struct {
	CinemaScreenID int
	Type           models.SeatType
}

type ListSeatsOutput struct {
	Seats []models.CinemaSeat `json:"seats"`
}

type GetSeatOutput struct {
	Seat models.CinemaSeat `json:"seat"`
}
