package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("BACKEND_SECRET_KEY"))

type claims struct {
	username string
	jwt.RegisteredClaims
}

type User struct {
	Username string `json:"user"`
}

func New(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(120 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Verify(tokenString string) error {
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

func Parse(tokenString string) (User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return User{}, err
	}

	parsedToken, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return User{}, fmt.Errorf("invalid token")
	}

	return User{Username: parsedToken.username}, nil
}
