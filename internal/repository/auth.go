package repository

import (
	"errors"

	"github.com/anzhy11/go-e-commerce/internal/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) GetRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := r.db.Where("token = ?", token).First(&refreshToken).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *AuthRepository) CreateRefreshToken(data *models.RefreshToken) error {
	return r.db.Create(&data).Error
}

func (r *AuthRepository) DeleteRefreshToken(data *models.RefreshToken) {
	r.db.Delete(&data)
}
