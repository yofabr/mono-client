package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/yofabr/mono-client/cmd/api"
	"github.com/yofabr/mono-client/cmd/application"
)

func validateEnvVars() error {
	required := []string{"PG_DSN", "REDIS_ADD", "REDIS_PASS", "REDIS_DB", "PORT", "SECRET"}
	for _, env := range required {
		if os.Getenv(env) == "" {
			return fmt.Errorf("missing required environment variable: %s", env)
		}
	}
	return nil
}

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Unable to load environmental variables")
	}

	if err := validateEnvVars(); err != nil {
		panic(err)
	}

	app := application.NewApp()
	app.Init()

	api := api.NewApi(*app)
	api.Init()

	port := ":" + os.Getenv("PORT")

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Error while starting the app:", err)
	}
}
