package user

import (
	"fmt"

	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/pkgs/models"
)

type CreateUserInput struct {
	Name     string          `json:"name"`
	Email    string          `json:"email"`
	Password string          `json:"password"`
	UserType models.UserType `json:"user_type"`
}

func (c *CreateUserInput) Validate(db *gorm.DB) error {
	if c.Name == "" || c.Email == "" || c.Password == "" {
		return fmt.Errorf("name, email, and password are required")
	}
	if !isValidUserType(c.UserType) {
		return fmt.Errorf("user_type must be one of: REGULAR, CINEMA_OWNER, ADMIN")
	}
	var existing models.User
	if err := db.Where("email = ?", c.Email).First(&existing).Error; err == nil {
		return fmt.Errorf("email already registered")
	}
	return nil
}

type CreateUserOutput struct {
	User models.User `json:"user"`
}

type ListUsersInput struct {
	Name     string
	Email    string
	UserType string
}

type ListUsersOutput struct {
	Users []models.User `json:"users"`
}

type GetUserByIDOutput struct {
	User models.User `json:"user"`
}

type UpdateUserInput struct {
	ID       int
	Name     string          `json:"name"`
	Email    string          `json:"email"`
	UserType models.UserType `json:"user_type"`
}

func (u *UpdateUserInput) Validate(db *gorm.DB) error {
	if u.UserType != "" && !isValidUserType(u.UserType) {
		return fmt.Errorf("user_type must be one of: REGULAR, CINEMA_OWNER, ADMIN")
	}
	return nil
}

type UpdateUserOutput struct {
	User models.User `json:"user"`
}

func isValidUserType(t models.UserType) bool {
	switch t {
	case models.UserTypeRegular, models.UserTypeCinemaOwner, models.UserTypeAdmin:
		return true
	}
	return false
}
