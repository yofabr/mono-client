package auth

import (
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yofabr/mono-client/cmd/application"
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

func (auth *AuthHandler) Register(IP string, email string, pass string) {
	// - check if the user exists in the db with this email.
	// - check if this client with IP address has other multiple accounts with IP
	// - check if credentials are unique and valid..
}

func (auth *AuthHandler) Login(IP string, email string, pass string) {
	// - check if the user has other accounts with this IP
	// - close other authenticated clients if there exist with this info
	// - generate secure jwt and set to IP list in redis..
}

func (auth *AuthHandler) GenerateToken(userID string, IP string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	secret := os.Getenv("SECRET")
	rd_client := auth.app.Databases.Redis()

	rd_client.Set(ctx, userID, IP, time.Hour*24*30)

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
