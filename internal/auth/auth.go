package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
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

// type AuthData struct {
// 	UserID string `json:"user_id"`
// 	IP     string `json:"ip"`
// 	Token  string `json:"token"`
// }

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

	token, err := auth.GenerateToken(userID)

	if err != nil {
		return "", err
	}

	if err := auth.SaveAuthentication(ctx, IP, userID, token); err != nil {
		return "", err
	}

	return token, nil
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
		return "", errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(userID)

	if err != nil {
		return "", err
	}

	if err := auth.CheckAuthAbility(ctx, IP, userID); err != nil {
		return "", err
	}

	if err := auth.SaveAuthentication(ctx, IP, userID, token); err != nil {
		return "", err
	}

	return token, nil
}

func (auth *AuthHandler) CheckAuthAbility(ctx context.Context, IP string, userId string) error {
	redis_client := auth.app.Databases.Redis()

	userIdKey := utils.KeyCacheUserID(userId)

	data, err := redis_client.Get(ctx, userIdKey).Result()
	if err != nil {
		return nil
	}

	data = strings.Split(data, "|-|")[1]

	fmt.Println("Existing IP:", data)
	fmt.Println("User IP:", IP)

	if data != IP {
		return errors.New("another device has already signed into this account")
	}

	// fmt.Println("User exists with userID:", UserIDAuthExists)

	//
	return nil
}

func (auth *AuthHandler) SaveAuthentication(ctx context.Context, IP, userId, token string) error {
	redisClient := auth.app.Databases.Redis()

	// data := AuthData{
	// 	UserID: userId,
	// 	IP:     IP,
	// 	Token:  token,
	// }

	// jsonData, err := json.Marshal(data)
	// if err != nil {
	// 	return err
	// }

	// Use pipeline to set all keys atomically
	pipe := redisClient.TxPipeline()
	UserIDKey := utils.KeyCacheUserID(userId)
	value := fmt.Sprintf("%v|-|%v", token, IP)

	pipe.Set(ctx, UserIDKey, value, 24*time.Hour)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return errors.New("timeout reached")
	default:
	}

	return nil
}

func (auth *AuthHandler) GenerateToken(userID string) (string, error) {
	secret := []byte(os.Getenv("SECRET"))

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
