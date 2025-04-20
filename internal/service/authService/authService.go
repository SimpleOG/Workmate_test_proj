package authService

import (
	"Workmate/internal/hashing"
	"Workmate/internal/jwt"
	db "Workmate/internal/repositories/postgresql/sqlc"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type AuthServiceInterface interface {
	Register(userParams db.CreateUserParams) (db.User, error)
	Login(userParams db.GetUserForLoginParams) (string, error)
	ValidateToken(tokenString string) (int32, error)
}

type AuthService struct {
	Querier    db.Querier
	Hasher     hashing.PassHasherInterface
	TokenMaker jwt.TokenMakerInterface
}

func NewAuthService(querier db.Querier, secretKey string) AuthServiceInterface {
	return &AuthService{
		Querier:    querier,
		Hasher:     hashing.NewPassHasher(),
		TokenMaker: jwt.NewTokenMaker(secretKey),
	}
}

func (a *AuthService) Register(userParams db.CreateUserParams) (db.User, error) {

	user, err := a.Querier.GetUserByUsername(context.Background(), userParams.Username)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, fmt.Errorf("cannot create new user %v", err)
		}
	}
	if user.ID != 0 {
		return db.User{}, errors.New("user already exists")
	}
	//создание пользователя
	hashedPassword, err := a.Hasher.HashPass(userParams.Password)
	if err != nil {
		return db.User{}, err
	}
	userParams.Password = hashedPassword
	NewUser, err := a.Querier.CreateUser(context.Background(), userParams)
	if err != nil {
		return db.User{}, fmt.Errorf("failed to create user %v", err)
	}
	return NewUser, nil
}
func (a *AuthService) Login(userParams db.GetUserForLoginParams) (string, error) {
	user, err := a.Querier.GetUserByUsername(context.Background(), userParams.Username)
	if err != nil {
		return "", fmt.Errorf("cannot find user by username %v", err)
	}
	if err := a.Hasher.ComparePass(userParams.Password, user.Password); err != nil {
		return "", fmt.Errorf("passwords don`t match %v", err)
	}
	token, err := a.TokenMaker.CreateToken(user.ID, time.Hour*24)
	if err != nil {
		return "", err
	}
	return token, nil
}
func (a *AuthService) ValidateToken(tokenString string) (int32, error) {
	claims, err := a.TokenMaker.ParseToken(tokenString)
	if err != nil {
		return 0, fmt.Errorf("Cannot validate token %v", err)
	}
	id, ok := claims["sub"].(float64)
	if !ok {
		return 0, fmt.Errorf("cannot convert claims %v", err)
	}
	user, err := a.Querier.GetUserByID(context.Background(), int32(id))
	if err != nil {
		return 0, fmt.Errorf("cannot find user %v", err)
	}
	return user.ID, nil
}
