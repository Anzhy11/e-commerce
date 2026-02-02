package repository

import (
	"errors"

	"github.com/anzhy11/go-e-commerce/internal/models"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(userID uint) (*models.User, error)
	CreateUser(data *models.User) error
	UpdateUser(data *models.User) error
}

type UserRpository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepositoryInterface {
	return &UserRpository{
		db: db,
	}
}

func (r *UserRpository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &user, nil
}

func (r *UserRpository) GetUserById(userID uint) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &user, nil
}

func (r *UserRpository) CreateUser(data *models.User) error {
	return r.db.Create(&data).Error
}

func (r *UserRpository) UpdateUser(data *models.User) error {
	return r.db.Save(&data).Error
}
