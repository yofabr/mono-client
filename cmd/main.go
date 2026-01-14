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
	err := godotenv.Load()

	if err != nil {
		panic("Unable to load environmental variables")
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
