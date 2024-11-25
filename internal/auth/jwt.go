package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minimumSecretLength = 32

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("the token has expired")
)

type JWTMaker struct {
	secret string
}

func NewJWTMaker(secret string) (Maker, error) {
	if len(secret) < minimumSecretLength {
		return nil, fmt.Errorf("invalid secret length , it must be atleast 32 characters")
	}

	return &JWTMaker{
		secret,
	}, nil
}

// This method is used to create a new JWT token , it implements the Maker interface
func (maker *JWTMaker) Create(email string, userId int64, duration time.Duration) (string, error) {
	payload := NewPayload(email, userId, duration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(maker.secret))
}

// This method is used to verify a new JWT token , it implements the Maker interface

func (maker *JWTMaker) Verify(tokenString string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(maker.secret), nil
	}
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, keyFunc)
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Payload)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// check if the token has expired
	if time.Now().After(claims.ExpiresAt) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
