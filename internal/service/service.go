package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shared-drawboard/internal/database"
	"github.com/shared-drawboard/internal/models"
	"github.com/shared-drawboard/pkg/auth"
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

func (s *Service) CreateSession(ctx context.Context, userID string) (*models.SessionDTO, error) {
	refreshtokenString, err := auth.CreateRefreshToken(32)
	if err != nil {
		return &models.SessionDTO{}, err
	}

	refreshtokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshtokenString), bcrypt.DefaultCost)
	if err != nil {
		return &models.SessionDTO{}, err
	}

	sDTO := models.SessionDTO{
		UserID:     userID,
		TokenHash:  string(refreshtokenHash),
		ExpiresAt:  strconv.FormatInt(time.Now().Add(24*7*time.Hour).Unix(), 10),
		CreatedAt:  strconv.FormatInt(time.Now().Unix(), 10),
		LastUsedAt: strconv.FormatInt(time.Now().Unix(), 10),
	}

	sDTO.ID, err = s.DB.CreateSession(ctx, sDTO)
	if err != nil {
		return &models.SessionDTO{}, err
	}

	return &sDTO, nil
}

func (s *Service) UpdateSession(ctx context.Context, tokenDTO models.RefreshTokenDTO) (*models.RefreshTokenDTO, error) {
	authToken := tokenDTO.AuthToken

	//verify jwt token
	claims, err := auth.VerifyJWTToken(authToken)
	if err != nil {
		return &models.RefreshTokenDTO{}, err
	}

	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return &models.RefreshTokenDTO{}, errors.New("invalid claims type")
	}

	uid := mapClaims["sub"].(string)

	//create new refresh token, auth token
	newAuthToken, err := auth.CreateJWTToken(uid, time.Now().Add(15*time.Minute).Unix())
	if err != nil {
		return &models.RefreshTokenDTO{}, err
	}

	newRefreshToken, err := auth.CreateRefreshToken(32)
	if err != nil {
		return &models.RefreshTokenDTO{}, err
	}

	newRefreshtokenHash, err := bcrypt.GenerateFromPassword([]byte(newRefreshToken), bcrypt.DefaultCost)
	if err != nil {
		return &models.RefreshTokenDTO{}, err
	}

	newTokenDTO := models.RefreshTokenDTO{
		AuthToken:        newAuthToken,
		AuthExpiresAt:    strconv.FormatInt(time.Now().Add(15*time.Minute).Unix(), 10),
		RefreshToken:     newRefreshToken,
		RefreshExpriesAt: strconv.FormatInt(time.Now().Add(24*7*time.Hour).Unix(), 10),
	}

	_, err = s.DB.UpdateSession(ctx, uid, string(newRefreshtokenHash))
	if err != nil {
		return &models.RefreshTokenDTO{}, err
	}

	return &newTokenDTO, nil
}
