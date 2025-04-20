package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type TokenMakerInterface interface {
	CreateToken(user_id int32, duration time.Duration) (string, error)
	ParseToken(tokenString string) (jwt.MapClaims, error)
}

type TokenMaker struct {
	secretKey string
}

func NewTokenMaker(secretKey string) TokenMakerInterface {
	return &TokenMaker{secretKey: secretKey}
}
func (t *TokenMaker) CreateToken(userId int32, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(duration).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", fmt.Errorf("cannot sign token %v", err)
	}
	return tokenStr, nil

}

func (t *TokenMaker) ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(t.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, errors.New("time limit expired")
		}
		return claims, nil
	} else {
		return nil, errors.New("cannot parse claims")
	}
}
