package main

import (
	"farstu/internal/config"
	"farstu/internal/views"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

func main() {
	// Config
	appConfig, err := config.ReadAppConfig("app.toml")
	if err != nil {
		log.Fatal(err)
	}

	// Routes
	http.Handle("/", templ.Handler(views.Hello("World")))

	// Run
	port := ":" + strconv.Itoa(appConfig.App.Port)
	log.Printf("App listening at %s", port)
	http.ListenAndServe(port, nil)
}
