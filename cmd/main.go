package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/yofabr/mono-client/cmd/api"
	"github.com/yofabr/mono-client/cmd/application"
)

func main() {
	// Load environment variables from .env so local development
	// behaves the same as deployment configuration.
	err := godotenv.Load()

	if err != nil {
		panic("Unable to load environmental variables")
	}

	app := application.NewApp()
	// Initialize shared infrastructure dependencies (Postgres and Redis).
	app.Init()

	api := api.NewApi(*app)
	// Register HTTP handlers on the default mux.
	api.Init()

	port := ":" + os.Getenv("PORT")

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Error while starting the app:", err)
	}
}
