package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/saleh-ghazimoradi/Projectopher/config"
	"time"
)

type Claims struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	UserId    string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(cfg *config.Config, firstname, lastname, email, role, userId string) (accessToken, RefreshToken string, err error) {
	accessClaims := &Claims{
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = at.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	refreshClaims := &Claims{
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.JWT.RefreshTokenExpires)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := rt.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
