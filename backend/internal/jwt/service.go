package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	key        []byte
	expiration time.Duration
}

func NewService(key []byte, expiration time.Duration) *Service {
	return &Service{
		key:        key,
		expiration: expiration,
	}
}

type claims struct {
	username string
	jwt.RegisteredClaims
}

type User struct {
	Username string `json:"user"`
}

func (s Service) New(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString(s.key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s Service) Verify(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func (s Service) Parse(tokenString string) (User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
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
