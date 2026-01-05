package main

import (
	"fmt"
	"net/http"

	"github.com/yofabr/mono-client/cmd/api"
	"github.com/yofabr/mono-client/cmd/application"
)

func main() {

	app := application.NewApp()
	app.Init()

	api := api.NewApi(*app)
	api.Init()

	err := http.ListenAndServe(":8080", nil)
	fmt.Println("Error while starting the app..:", err)
}
