package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/yofabr/mono-client/cmd/application"
	"github.com/yofabr/mono-client/internal/utils"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandler struct {
	app *application.Application
}

func NewAuthHandler(app *application.Application) *AuthHandler {
	return &AuthHandler{app: app}
}

func (auth *AuthHandler) Register(IP, email, pass string) (string, error) {
	db := auth.app.Databases.PG()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hashedPW, err := utils.Hash(pass)
	if err != nil {
		return "", err
	}

	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id;
	`

	var userID string
	err = db.QueryRow(ctx, query, email, hashedPW).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (auth *AuthHandler) Login(IP, email, pass string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := auth.app.Databases.PG()

	query := `
		SELECT id, password
		FROM users
		WHERE email = $1
	`

	var (
		userID       string
		passwordHash string
	)

	err := db.QueryRow(ctx, query, email).Scan(&userID, &passwordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	if err := utils.Compare(pass, passwordHash); err != nil {
		return "", errors.New("Invalid credentials")
	}

	return userID, nil
}

func (auth *AuthHandler) GenerateToken(userID string, IP string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	secret := os.Getenv("SECRET")
	rd_client := auth.app.Databases.Redis()

	rd_client.Set(ctx, userID, IP, time.Hour*24*30)

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
