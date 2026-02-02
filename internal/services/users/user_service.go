package userService

import (
	"errors"

	"github.com/anzhy11/go-e-commerce/internal/dto"
	"github.com/anzhy11/go-e-commerce/internal/repository"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	GetProfile(userID uint) (*dto.UserResponse, error)
	UpdateProfile(userID uint, data *dto.UpdateProfileRequest) (*dto.UserResponse, error)
}

type userService struct {
	db       *gorm.DB
	userRepo repository.UserRepositoryInterface
}

func New(db *gorm.DB) UserServiceInterface {
	return &userService{
		db:       db,
		userRepo: repository.NewUserRepo(db),
	}
}

func (s *userService) GetProfile(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      user.Role,
		IsActive:  user.IsActive,
	}, nil
}

func (s *userService) UpdateProfile(userID uint, data *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	user.FirstName = data.FirstName
	user.LastName = data.LastName
	user.Phone = data.Phone

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      user.Role,
		IsActive:  user.IsActive,
	}, nil
}
