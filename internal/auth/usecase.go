package auth

import (
	"errors"
	"time"

	env "github/com/cl0ky/e-voting-be/env"
	"github/com/cl0ky/e-voting-be/models"
	jwtutil "github/com/cl0ky/e-voting-be/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UseCase interface {
	Register(req RegisterRequest) (RegisterResponse, error)
	Login(req LoginRequest) (LoginResponse, error)
	GetProfile(userId string) (*models.User, error)
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
		return RegisterResponse{}, errors.New("Email atau NIK sudah terdaftar")
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
	user.Id = uuid.New()

	user.BaseModel.CreatedBy = &user.Id

	if err := us.repo.CreateUser(&user); err != nil {
		return RegisterResponse{}, errors.New("gagal register user")
	}

	return RegisterResponse{
		Message: "Registrasi berhasil",
	}, nil
}

func (us *useCase) Login(req LoginRequest) (LoginResponse, error) {
	user, err := us.repo.GetUserByEmailOrNIK(req.EmailOrNIK)
	if err != nil {
		return LoginResponse{}, errors.New("email/nik atau password salah")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return LoginResponse{}, errors.New("email/nik atau password salah")
	}

	jwtManager := jwtutil.NewJWTManager(env.JWTSecret, 24*time.Hour)
	token, err := jwtManager.Generate(user.Id.String(), user.Role)
	if err != nil {
		return LoginResponse{}, errors.New("gagal generate token")
	}

	return LoginResponse{
		Token:   token,
		Role:    user.Role,
		Message: "Login berhasil",
	}, nil
}

func (us *useCase) GetProfile(userId string) (*models.User, error) {
	return us.repo.GetUserByID(userId)
}
