package authService

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/anzhy11/go-e-commerce/internal/dto"
	"github.com/anzhy11/go-e-commerce/internal/models"
	"github.com/anzhy11/go-e-commerce/internal/repository"
	"github.com/anzhy11/go-e-commerce/internal/utils"
	"github.com/anzhy11/go-e-commerce/pkg/encryption"
	"github.com/rs/zerolog"
)

type AuthServiceInterface interface {
	Register(data *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(data *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(data *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(rt string) error
}

type authService struct {
	db        *gorm.DB
	log       *zerolog.Logger
	cfg       *config.Config
	userRepo  *repository.UserRpository
	orderRepo *repository.OrderRepository
	authRepo  *repository.AuthRepository
}

func New(db *gorm.DB, cfg *config.Config, log *zerolog.Logger) AuthServiceInterface {
	return &authService{
		db:        db,
		cfg:       cfg,
		log:       log,
		userRepo:  repository.NewUserRepo(db),
		orderRepo: repository.NewOrderRepo(db),
		authRepo:  repository.NewAuthRepo(db),
	}
}

func (s *authService) Register(data *dto.RegisterRequest) (*dto.AuthResponse, error) {
	existingUser, err := s.userRepo.GetUserByEmail(data.Email)
	if err != nil {
		return nil, err
	}

	if existingUser.ID != 0 {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := encryption.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Password:  hashedPassword,
		Phone:     data.Phone,
		Role:      string(models.RoleCustomer),
	}

	if err := s.userRepo.CreateUser(&user); err != nil {
		return nil, err
	}

	cart := models.Cart{
		UserID: user.ID,
	}

	if err := s.orderRepo.CreateCart(&cart); err != nil {
		s.log.Error().Err(err).Msg("Failed to create cart")
	}

	return s.generateAuthResponse(&user)
}

func (s *authService) Login(data *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetUserByEmail(data.Email)
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	if !encryption.CheckPassword(data.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	return s.generateAuthResponse(user)
}

func (s *authService) RefreshToken(data *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	payload, err := utils.VerifyToken(data.RefreshToken, s.cfg.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	refreshToken, err := s.authRepo.GetRefreshToken(data.RefreshToken)
	if err != nil {
		return nil, err
	}

	if refreshToken.UserID != payload.UserID {
		return nil, errors.New("invalid refresh token")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	user, err := s.userRepo.GetUserById(payload.UserID)
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	s.authRepo.DeleteRefreshToken(refreshToken)

	return s.generateAuthResponse(user)
}

func (s *authService) Logout(rt string) error {
	return s.db.Where("token = ?", rt).Delete(&models.RefreshToken{}).Error
}

func (s *authService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(s.cfg, user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	refreshTokenModel := models.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(s.cfg.JWT.RefreshTokenExpiresIn),
	}

	if err := s.authRepo.CreateRefreshToken(&refreshTokenModel); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
