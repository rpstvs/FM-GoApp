package main

import (
	"net/http"

	"github.com/rpstvs/fm-goapp/internal/app"
)

func main() {
	app, err := app.NewApplication()

	if err != nil {
		panic(err)
	}

	app.Logger.Println("we Are running!")

	server := &http.Server{
		Addr:        ":8080",
		IdleTimeout: 30 * Time.Second,
	}
}
