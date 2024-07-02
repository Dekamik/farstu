package main

import (
	"farstu/internal/config"
	"farstu/internal/views"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/a-h/templ"
)

func main() {
	// Config
	appConfigPath := "app.toml"
	appConfig, err := config.ReadAppConfig(appConfigPath)
	if err != nil {
		slog.Error("An error ocurred while reading " + appConfigPath, 
			"err", err,
		)
		os.Exit(1)
	}

	// Routes
	http.Handle("/", templ.Handler(views.Hello("World")))

	// Run
	port := ":" + strconv.Itoa(appConfig.App.Port)
	slog.Info("App started", "port", port)
	http.ListenAndServe(port, nil)
}
