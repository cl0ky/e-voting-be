package auth

import (
	"github/com/cl0ky/e-voting-be/models"

	"gorm.io/gorm"
)

type Repository interface {
	IsEmailOrNIKExist(email, nik string) (bool, error)
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByEmailOrNIK(emailOrNIK string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) IsEmailOrNIKExist(email, nik string) (bool, error) {
	var user models.User
	err := r.db.Where("email = ? OR nik = ?", email, nik).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *repository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *repository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetUserByEmailOrNIK(emailOrNIK string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? OR nik = ?", emailOrNIK, emailOrNIK).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
