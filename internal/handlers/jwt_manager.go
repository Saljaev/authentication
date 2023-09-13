package handlers

import (
	"authentication/internal/entity"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"time"
)

type JWTManager struct {
	secret        string
	tokenDuration time.Duration
}

type UseClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

func NewJWTManager(secret string, tokenDuratino time.Duration) *JWTManager {
	return &JWTManager{secret, tokenDuratino}
}

func (m *JWTManager) Generate(user *entity.User) (string, error) {
	claims := UseClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenDuration)),
		},
		Email: user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(m.secret))
}

func (m *JWTManager) Parse(accessToken string) (*UseClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UseClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return []byte(m.secret), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UseClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (m *JWTManager) NewRefreshToken() (string, error) {
	tokenLength, err := strconv.ParseInt(os.Getenv("token_length"), 10, 64)
	if err != nil {
		return "", fmt.Errorf("Error to convert token length to int: %w", err)
	}
	token := make([]byte, tokenLength)
	_, err = rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("JWTManager - NewRefreshToken - Generate")
	}

	return base64.StdEncoding.EncodeToString(token), nil
}
