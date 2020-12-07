package main

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTClaims struct {
	UserId int `json:"u"`
	jwt.StandardClaims
}

const JWTDuration = time.Hour * 24 * 7

func SignJWT(secret string, userId int) (string, error) {
	if secret == "" {
		return "", errors.New("JWTSecret in config is blank")
	}
	claims := JWTClaims{
		UserId: userId,
	}
	claims.ExpiresAt = time.Now().Add(JWTDuration).UnixNano()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}
