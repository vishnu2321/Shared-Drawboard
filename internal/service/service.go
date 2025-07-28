package service

import (
	"context"
	"errors"

	"github.com/shared-drawboard/internal/database"
	"github.com/shared-drawboard/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	DB database.DB
}

func New() (s *Service, err error) {
	db, err := database.New()
	if err != nil {
		return nil, err
	}
	return &Service{DB: db}, nil
}

var ErrUserExists = errors.New("user already exists")

func (s *Service) SaveUser(ctx context.Context, u models.User) (id string, err error) {
	existing, err := s.DB.FindBy(ctx, "Email", u.Email)
	if err != nil {
		return "", err
	}
	if existing != nil {
		return "", ErrUserExists
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	u.Password = string(hashedpassword)

	return s.DB.SaveUserDB(ctx, u)
}

func (s *Service) GetUser(ctx context.Context, value string) (*models.User, error) {
	user, err := s.DB.FindBy(ctx, "Email", value)
	if err != nil {
		return &models.User{}, err
	}
	u := models.User{
		ID:       string(user.ID.Hex()),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	return &u, nil
}
