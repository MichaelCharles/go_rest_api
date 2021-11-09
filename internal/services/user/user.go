package user

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/mcaubrey/go_rest_api/internal/services/comment"
)

// Service - base service struct
type Service struct {
	DB *gorm.DB
}

// NewService - returns a new service
func NewService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}

// User - model for users
type User struct {
	gorm.Model
	Email      string
	Username   string
	Password   string
	GivenName  string
	FamilyName string
	Comments   []comment.Comment `gorm:"many2many:user_comment;"`
}

// UserService - interface for user service
type UserService interface {
	GetUser(ID uint) (User, error)
	LoginUser(identifier string, password string) (User, error)
	RegisterUser(newUser User) (User, error)
}

// GetUser - gets user by ID from database
func (s *Service) GetUser(ID uint) (User, error) {
	var user User
	if result := s.DB.First(&user, ID); result.Error != nil {
		return User{}, result.Error
	}
	return user, nil
}

// LoginUser - signs a user in and returns a JWT
func (s *Service) LoginUser(credentials User) (User, error) {
	var user User
	if result := s.DB.First(&user, "email = ?", credentials.Email); result.Error != nil {
		if result := s.DB.First(&user, "username = ?", credentials.Username); result.Error != nil {
			return User{}, result.Error
		}
	}

	if user.Password != credentials.Password {
		return User{}, errors.New("incorrect password")
	}
	return user, nil
}

// RegisterUser - adds a new user to the database
func (s *Service) RegisterUser(newUser User) (User, error) {
	if result := s.DB.Save(&newUser); result.Error != nil {
		return User{}, result.Error
	}
	return newUser, nil
}

// GetAllUsers - retrieves all users without restrictions
func (s *Service) GetAllUsers() ([]User, error) {
	var users []User
	if result := s.DB.Find(&users); result.Error != nil {
		return users, result.Error
	}
	return users, nil
}
