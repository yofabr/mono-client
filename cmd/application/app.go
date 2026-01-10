package application

import (
	"log"
	"os"
	"strconv"

	"github.com/yofabr/mono-client/internal/databases"
)

type Application struct {
	Databases *databases.Databases
}

func NewApp() *Application {
	db := databases.Databases{}

	app := Application{}

	app.Databases = &db

	return &app
}

func (a *Application) Init() {
	// godotenv.Load()

	defer log.Println("Application is initialized successfully")

	pg_dns := os.Getenv("PG_DSN")
	redis_add := os.Getenv("REDIS_ADD")
	redis_pass := os.Getenv("REDIS_PASS")
	redis_db, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		log.Fatal("Error while parsing redis db")
	}
	a.Databases.NewPostgresInit(pg_dns)                   // Postgres init
	a.Databases.NewRedis(redis_add, redis_pass, redis_db) /// Redis init
}
