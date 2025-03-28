package jwt

import (
	"backend/internal"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	logger     *zap.Logger
	secret     string
	expiration time.Duration
}

func NewService(logger *zap.Logger, secret string, expiration time.Duration) *Service {
	return &Service{
		logger:     logger,
		secret:     secret,
		expiration: expiration,
	}
}

type claims struct {
	id       string
	username string
	role     string
	jwt.RegisteredClaims
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"user"`
	Role     string `json:"role"`
}

func (s Service) New(ctx context.Context, id, username string, role string) (string, error) {
	logger := internal.LoggerWithContext(ctx, s.logger)

	jwtID := uuid.New()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		id:       id,
		username: username,
		role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "CLI-Forum",
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jwtID.String(),
		},
	})

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		logger.Error("Failed to sign token", zap.Error(err), zap.String("id", id), zap.String("username", username), zap.String("role", role))
		return "", err
	}

	logger.Debug("Generated new JWT token", zap.String("id", id), zap.String("username", username), zap.String("role", role))

	return tokenString, nil
}

func (s Service) Parse(ctx context.Context, tokenString string) (User, error) {
	logger := internal.LoggerWithContext(ctx, s.logger)

	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			logger.Warn("Failed to parse JWT token due to malformed structure, this is not a JWT token", zap.String("error", err.Error()))
			return User{}, err
		case errors.Is(err, jwt.ErrSignatureInvalid):
			logger.Warn("Failed to parse JWT token due to invalid signature", zap.String("error", err.Error()))
			return User{}, err
		case errors.Is(err, jwt.ErrTokenExpired):
			expiredTime, getErr := token.Claims.GetExpirationTime()
			if getErr != nil {
				logger.Warn("Failed to parse JWT token due to expired timestamp", zap.String("error", err.Error()))
			} else {
				logger.Warn("Failed to parse JWT token due to expired timestamp", zap.String("error", err.Error()), zap.Time("expired_at", expiredTime.Time))
			}

			return User{}, err
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			notBeforeTime, getErr := token.Claims.GetNotBefore()
			if getErr != nil {
				logger.Warn("Failed to parse JWT token due to not valid yet timestamp", zap.String("error", err.Error()))
			} else {
				logger.Warn("Failed to parse JWT token due to not valid yet timestamp", zap.String("error", err.Error()), zap.Time("not_valid_yet", notBeforeTime.Time))
			}

			return User{}, err
		default:
			logger.Error("Failed to parse or validate JWT token", zap.Error(err))
		}
	}

	claims, ok := token.Claims.(*claims)
	if !ok {
		logger.Error("Failed to extract claims from JWT token")
		return User{}, fmt.Errorf("failed to extract claims from JWT token")
	}

	logger.Debug("Successfully parsed JWT token", zap.String("id", claims.id), zap.String("username", claims.username), zap.String("role", claims.role))

	return User{
		ID:       claims.id,
		Username: claims.username,
		Role:     claims.role,
	}, nil
}
