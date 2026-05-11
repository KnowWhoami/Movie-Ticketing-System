package user

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"gorm.io/gorm"

	"KnowWhoami/movie-ticketing/pkgs/auth"
	"KnowWhoami/movie-ticketing/pkgs/models"
)

type Service struct {
	db     *gorm.DB
	logger log.Logger
}

func NewService(db *gorm.DB, logger log.Logger) *Service {
	return &Service{db: db, logger: logger}
}

func (s *Service) CreateUser(inp *CreateUserInput) (*CreateUserOutput, error) {
	u := models.User{
		Name:     inp.Name,
		Email:    inp.Email,
		Password: auth.HashPassword(inp.Password),
		UserType: inp.UserType,
	}
	if err := s.db.Create(&u).Error; err != nil {
		_ = s.logger.Log("method", "CreateUser", "err", err)
		return nil, err
	}
	_ = s.logger.Log("method", "CreateUser", "user_id", u.ID)
	return &CreateUserOutput{User: u}, nil
}

func (s *Service) ListUsers(inp *ListUsersInput) (*ListUsersOutput, error) {
	query := s.db.Model(&models.User{})
	if inp.Name != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+inp.Name+"%")
	}
	if inp.Email != "" {
		query = query.Where("LOWER(email) LIKE LOWER(?)", "%"+inp.Email+"%")
	}
	if inp.UserType != "" {
		query = query.Where("user_type = ?", inp.UserType)
	}
	var users []models.User
	if err := query.Find(&users).Error; err != nil {
		_ = s.logger.Log("method", "ListUsers", "err", err)
		return nil, err
	}
	return &ListUsersOutput{Users: users}, nil
}

func (s *Service) GetUserByID(id int) (*GetUserByIDOutput, error) {
	var u models.User
	if err := s.db.First(&u, id).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &GetUserByIDOutput{User: u}, nil
}

func (s *Service) UpdateUser(inp *UpdateUserInput) (*UpdateUserOutput, error) {
	var u models.User
	if err := s.db.First(&u, inp.ID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	if inp.Email != "" && inp.Email != u.Email {
		var existing models.User
		if err := s.db.Where("email = ? AND id != ?", inp.Email, inp.ID).First(&existing).Error; err == nil {
			return nil, fmt.Errorf("email already in use")
		}
		u.Email = inp.Email
	}
	if inp.Name != "" {
		u.Name = inp.Name
	}
	if inp.UserType != "" {
		u.UserType = inp.UserType
	}
	if err := s.db.Save(&u).Error; err != nil {
		_ = s.logger.Log("method", "UpdateUser", "err", err)
		return nil, err
	}
	_ = s.logger.Log("method", "UpdateUser", "user_id", u.ID)
	return &UpdateUserOutput{User: u}, nil
}

func (s *Service) DeleteUser(id int) error {
	result := s.db.Delete(&models.User{}, id)
	if result.Error != nil {
		_ = s.logger.Log("method", "DeleteUser", "err", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	_ = s.logger.Log("method", "DeleteUser", "user_id", id)
	return nil
}
