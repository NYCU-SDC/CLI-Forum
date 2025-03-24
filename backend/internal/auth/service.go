package auth

import (
	"context"
	"errors"
	"os"

	"backend/internal/jwt"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func connect(ctx context.Context) (*pgx.Conn, error) {
	dbURL := os.Getenv("DATABASE_URL")
	return pgx.Connect(ctx, dbURL)
}

func CreateUser(ctx context.Context, u RegisterRequest) (User, error) {
	// Connect to the database
	conn, err := connect(ctx)
	if err != nil {
		return User{}, err
	}
	defer conn.Close(ctx)

	// Hash the password
	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return User{}, err
	}

	// Query the database
	queries := New(conn)
	return queries.Create(ctx, CreateParams{
		Name:     u.Username,
		Password: hashedPassword,
	})
}

func FindByName(ctx context.Context, username string) (User, error) {
	// Connect to the database
	conn, err := connect(ctx)
	if err != nil {
		return User{}, err
	}
	defer conn.Close(ctx)

	// Query the database
	queries := New(conn)

	return queries.FindByName(ctx, username)
}

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
