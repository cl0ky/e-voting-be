package auth

import (
	"errors"
	"github/com/cl0ky/e-voting-be/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UseCase interface {
	Register(req RegisterRequest) (RegisterResponse, error)
}

type useCase struct {
	repo Repository
	db   *gorm.DB
}

func NewUseCase(db *gorm.DB, repo Repository) UseCase {
	return &useCase{
		repo: repo,
		db:   db,
	}
}

func (us *useCase) Register(req RegisterRequest) (RegisterResponse, error) {
	exist, err := us.repo.IsEmailOrNIKExist(req.Email, req.NIK)
	if err != nil {
		return RegisterResponse{}, err
	}
	if exist {
		return RegisterResponse{}, errors.New("email atau NIK sudah terdaftar")
	}

	rtUUID, err := uuid.Parse(req.RTId)
	if err != nil {
		return RegisterResponse{}, errors.New("RT ID tidak valid")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return RegisterResponse{}, errors.New("gagal hash password")
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
		NIK:      req.NIK,
		RTId:     &rtUUID,
	}

	if err := us.repo.CreateUser(&user); err != nil {
		return RegisterResponse{}, errors.New("gagal register user")
	}

	return RegisterResponse{
		UserID:  user.Id.String(),
		Message: "Registrasi berhasil",
	}, nil
}
