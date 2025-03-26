package auth

import (
	"backend/internal/jwt"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *pgxpool.Pool

	count int
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

var hashCost = 14

func (s Service) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return string(bytes), err
}

func (s Service) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s Service) CreateUser(ctx context.Context, u User) (User, error) {
	// Get the connection form the pool
	queries := New(s.db)

	// Query the database
	return queries.Create(ctx, CreateParams{
		Name:     u.Name,
		Password: u.Password,
	})
}

func (s Service) FindUserByName(ctx context.Context, username string) (User, error) {
	// Get the connection form the pool
	queries := New(s.db)

	// Query the database
	return queries.FindByName(ctx, username)
}

func (s Service) IsUserExist(ctx context.Context, username string) (bool, error) {
	// Get the connection form the pool
	queries := New(s.db)

	// Query the database
	_, err := queries.FindByName(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s Service) Register(ctx context.Context, u RegisterRequest) (string, error) {
	// Check if the user already exists
	isUserExist, err := s.IsUserExist(ctx, u.Username)
	if err != nil {
		return "", errors.New("error_checking_user")
	}
	if isUserExist {
		return "", errors.New("user_already_exists")
	}

	// Create the user
	hashedPassword, err := s.HashPassword(u.Password)
	if err != nil {
		return "", errors.New("error_hashing_password")
	}

	user := User{
		Name:     u.Username,
		Password: hashedPassword,
	}

	_, err = s.CreateUser(ctx, user)
	if err != nil {
		fmt.Println("when creating user: ", err)
		return "", errors.New("error_registering_user")
	}

	// Registration successful, generate a token
	tokenString, err := jwt.New(u.Username)
	if err != nil {
		return "", errors.New("error_generating_token")
	}

	return tokenString, err
}

func (s Service) Login(ctx context.Context, u LoginRequest) (string, error) {
	// Find user by names
	user, err := s.FindUserByName(ctx, u.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("user_not_found")
		}
		return "", errors.New("error_finding_user")
	}

	// Check if the password is correct
	if !s.VerifyPassword(u.Password, user.Password) {
		return "", errors.New("incorrect_password")
	}

	// Login successful, generate a token
	tokenString, err := jwt.New(u.Username)
	if err != nil {
		return "", errors.New("error_generating_token")
	}

	return tokenString, err
}
