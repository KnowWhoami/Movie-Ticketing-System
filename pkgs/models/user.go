package models

type UserType string

const (
	UserTypeRegular     UserType = "REGULAR"
	UserTypeCinemaOwner UserType = "CINEMA_OWNER"
	UserTypeAdmin       UserType = "ADMIN"
)

type User struct {
	Model
	Name     string
	Email    string   `gorm:"type:varchar(191);uniqueIndex"`
	Password string   `json:"-"`
	UserType UserType `json:"user_type" gorm:"default:'REGULAR'"`

	// relations
	Bookings []Booking `gorm:"foreignKey:UserID"`
}
