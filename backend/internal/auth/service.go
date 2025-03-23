package auth

import (
	"errors"

	"backend/internal/jwt"
)

func Register(u RegisterRequest) (string, error) {
	// TODO: Implement registration with the database

	// Registration successful, generate a token
	tokenString, err := jwt.New(u.Username)
	if err != nil {
		return "", errors.New("error_generating_token")
	}

	return tokenString, err
}

func Login(u LoginRequest) (string, error) {
	// TODO: implement login with database

	// Registration successful, generate a token
	tokenString, err := jwt.New(u.Username)
	if err != nil {
		return "", errors.New("error_generating_token")
	}

	return tokenString, err
}
