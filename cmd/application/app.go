package application

import (
	"log"
	"os"
	"strconv"

	"github.com/yofabr/mono-client/internal/databases"
)

// Application is the root container for shared infrastructure services.
type Application struct {
	Databases *databases.Databases
}

// NewApp constructs the application shell with an initialized database wrapper.
func NewApp() *Application {
	db := databases.Databases{}

	app := Application{}
	app.Databases = &db

	return &app
}

// Init loads runtime configuration and initializes backing services.
func (a *Application) Init() {
	defer log.Println("Application is initialized successfully")

	// Environment variables are provided by .env (local) or deployment configs.
	pgDSN := os.Getenv("PG_DSN")
	redisAddr := os.Getenv("REDIS_ADD")
	redisPass := os.Getenv("REDIS_PASS")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal("Error while parsing redis db")
	}

	// Initialize persistent storage first, then cache/session store.
	a.Databases.NewPostgresInit(pgDSN)
	a.Databases.NewRedis(redisAddr, redisPass, redisDB)
}
