package models

type SeatType string

const (
	Recliner SeatType = "RECLINER"
	Premium  SeatType = "PREMIUM"
	FrontRow SeatType = "FRONT"
	Balcony  SeatType = "BALCONY"
)

type City struct {
	Model
	Name    string `json:"name"`
	ZipCode string `json:"zip_code"`
}

type Cinema struct {
	Model
	Name string `json:"name"`

	// relations
	CinemaOwnerID int            `json:"-"`
	CinemaOwner   User           `json:"cinema_owner"`
	CinemaScreens []CinemaScreen `json:"screens" gorm:"foreignkey:CinemaID"`
	CityID        int            `json:"-"`
	City          City           `json:"city"`
}

type CinemaScreen struct {
	Model
	Name string `json:"name" gorm:"uniqueIndex:unique_screen_per_cinema"`

	// relations
	CinemaID    int          `json:"-" gorm:"uniqueIndex:unique_screen_per_cinema"`
	Cinema      Cinema       `json:"-"`
	CinemaSeats []CinemaSeat `json:"seats" gorm:"foreignkey:CinemaScreenID"`
}

type CinemaSeat struct {
	Model
	SeatNumber int      `json:"seat_number" gorm:"uniqueIndex:unique_seat_per_screen"`
	Type       SeatType `json:"type"`

	// relations
	CinemaScreenID int          `json:"-" gorm:"uniqueIndex:unique_seat_per_screen"`
	CinemaScreen   CinemaScreen `json:"-"`
}
