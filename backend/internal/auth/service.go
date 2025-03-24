package auth

import (
	"context"
	"errors"
	"fmt"
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

func CreateUser(ctx context.Context, u User) (User, error) {
	// Connect to the database
	conn, err := connect(ctx)
	if err != nil {
		return User{}, err
	}
	defer conn.Close(ctx)

	// Query the database
	queries := New(conn)
	return queries.Create(ctx, CreateParams{
		Name:     u.Name,
		Password: u.Password,
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

func Exist(ctx context.Context, username string) (bool, error) {
	// Connect to the database
	conn, err := connect(ctx)
	if err != nil {
		return false, err
	}
	defer conn.Close(ctx)

	// Query the database
	queries := New(conn)

	_, err = queries.FindByName(ctx, username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func Register(u RegisterRequest) (string, error) {
	// Check if the user already exists
	isUserExist, err := Exist(context.Background(), u.Username)
	if err != nil {
		return "", errors.New("error_checking_user")
	}
	if isUserExist {
		return "", errors.New("user_already_exists")
	}

	// Create the user
	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return "", errors.New("error_hashing_password")
	}

	user := User{
		Name:     u.Username,
		Password: hashedPassword,
	}

	_, err = CreateUser(context.Background(), user)
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

func Login(u LoginRequest) (string, error) {
	// Find user by names
	user, err := FindByName(context.Background(), u.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", errors.New("user_not_found")
		}
		return "", errors.New("error_finding_user")
	}

	// Check if the password is correct
	if !VerifyPassword(u.Password, user.Password) {
		return "", errors.New("incorrect_password")
	}

	// Login successful, generate a token
	tokenString, err := jwt.New(u.Username)
	if err != nil {
		return "", errors.New("error_generating_token")
	}

	return tokenString, err
}
