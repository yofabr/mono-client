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
	godotenv.Load()

	app := application.NewApp()
	app.Init()

	api := api.NewApi(*app)
	api.Init()

	port := os.Getenv("PORT")

	err := http.ListenAndServe(port, nil)
	fmt.Println("Error while starting the app..:", err)
}
