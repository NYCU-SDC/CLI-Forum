package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("BACKEND_SECRET_KEY"))

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(120 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func Register(u RegisterRequest) (string, error) {
	// TODO: Implement registration with the database

	// Registration successful, generate a token
	tokenString, err := CreateToken(u.Username)
	if err != nil {
		return "", errors.New("error_generating_token")
	}

	return tokenString, err
}

func Login(u LoginRequest) (string, error) {
	// TODO: implement login with database

	// Registration successful, generate a token
	tokenString, err := CreateToken(u.Username)
	if err != nil {
		return "", errors.New("error_generating_token")
	}

	return tokenString, err
}
