package main

import (
	"farstu/internal/config"
	"farstu/internal/models"
	"farstu/internal/views"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/a-h/templ"
)

func handleTempl(path string, componentFunc func() templ.Component) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		componentFunc().Render(r.Context(), w)
	})
}

func main() {
	// Config
	appConfigPath := "app.toml"

	appConfig, err := config.ReadAppConfig(appConfigPath)
	if err != nil {
		slog.Error("An error ocurred while reading "+appConfigPath,
			"err", err,
		)
		os.Exit(1)
	}

	// Routes
	handleTempl("/", func() templ.Component { 
		return views.Index(models.NewIndexModel())
	})
	handleTempl("/htmx/time", func() templ.Component { 
		return views.Time(models.NewTimeModel()) 
	})

	// Static files
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.Handle("/static/", fs)

	// Run
	port := ":" + strconv.Itoa(appConfig.App.Port)
	slog.Info("Farstu started", "port", port)
	http.ListenAndServe(port, nil)
}
