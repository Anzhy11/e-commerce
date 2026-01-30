package utils

import (
	"errors"
	"time"

	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateTokenPair(cfg *config.Config, userID uint, email, role string) (accessToken, refreshToken string, err error) {
	accessPayload := &Payload{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessPayload)
	accessToken, err = at.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return
	}

	refreshPayload := &Payload{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.RefreshTokenExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshPayload)
	refreshToken, err = rt.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return
	}

	return accessToken, refreshToken, nil
}

func VerifyToken(tokenString, secret string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if payload, ok := token.Claims.(*Payload); ok && token.Valid {
		return payload, nil
	}

	return nil, errors.New("invalid token")
}
