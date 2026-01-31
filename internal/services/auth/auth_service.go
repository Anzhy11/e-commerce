package authService

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/anzhy11/go-e-commerce/internal/dto"
	"github.com/anzhy11/go-e-commerce/internal/models"
	"github.com/anzhy11/go-e-commerce/internal/utils"
	"github.com/anzhy11/go-e-commerce/pkg/encryption"
	"github.com/rs/zerolog"
)

type AuthServiceInterface interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(rt string) error
}

type AuthService struct {
	db  *gorm.DB
	log *zerolog.Logger
	cfg *config.Config
}

func New(db *gorm.DB, cfg *config.Config, log *zerolog.Logger) AuthServiceInterface {
	return &AuthService{
		db:  db,
		cfg: cfg,
		log: log,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existingUser.ID != 0 {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := encryption.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
		Phone:     req.Phone,
		Role:      string(models.RoleCustomer),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	cart := models.Cart{
		UserID: user.ID,
	}

	if err := s.db.Create(&cart).Error; err != nil {
		s.log.Error().Err(err).Msg("Failed to create cart")
	}

	return s.generateAuthResponse(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	if !encryption.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	return s.generateAuthResponse(&user)
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	payload, err := utils.VerifyToken(req.RefreshToken, s.cfg.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var refreshToken models.RefreshToken
	if err := s.db.Where("token = ?", req.RefreshToken).First(&refreshToken).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if refreshToken.UserID != payload.UserID {
		return nil, errors.New("invalid refresh token")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	var user models.User
	if err := s.db.Where("id = ?", payload.UserID).First(&user).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	s.db.Delete(&refreshToken)

	return s.generateAuthResponse(&user)
}

func (s *AuthService) Logout(rt string) error {
	return s.db.Where("token = ?", rt).Delete(&models.RefreshToken{}).Error
}

func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(s.cfg, user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	refreshTokenModel := models.RefreshToken{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(s.cfg.JWT.RefreshTokenExpiresIn),
	}

	if err := s.db.Create(&refreshTokenModel).Error; err != nil {
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
