package jwt

import (
	"errors"
	"go.uber.org/zap"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	logger     *zap.Logger
	key        []byte
	expiration time.Duration
}

func NewService(logger *zap.Logger, key []byte, expiration time.Duration) *Service {
	return &Service{
		logger:     logger,
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
		s.logger.Error("error when generating JWT token", zap.String("username", username), zap.Error(err))
		return "", err
	}

	return tokenString, nil
}

func (s Service) Verify(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})

	if err != nil {
		s.logger.Error("error when generating JWT token", zap.String("token", tokenString), zap.Error(err))
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
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
		return User{}, errors.New("invalid token")
	}

	return User{Username: parsedToken.username}, nil
}
