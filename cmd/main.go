package main

import (
	"farstu/internal/clock"
	"farstu/internal/config"
	"farstu/internal/index"
	"farstu/internal/shared"
	"farstu/internal/yr"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/a-h/templ"
)

var logLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

var navItems = []shared.NavItemModel{
	{
		Href: "/",
		Icon: "bi-house-door-fill",
	},
	{
		Href: "/weather",
		Icon: "bi-cloud-sun-fill",
	},
	{
		Href: "/disruptions",
		Icon: "bi-exclamation-triangle-fill",
	},
}

func getPageModel(activeHref string) shared.PageModel {
	navModels := make([]shared.NavItemModel, 0)

	for _, item := range navItems {
		model := shared.NavItemModel{
			Href: item.Href,
			Icon: item.Icon,
			IsActive: string(item.Href) == activeHref,
		}
		navModels = append(navModels, model)
	}

	return shared.PageModel{
		NavItems: navModels,
	}
}

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
		slog.Error("An error occurred while reading "+appConfigPath,
			"err", err,
		)
		os.Exit(1)
	}

	// Logging
	logLevel := logLevelMap[appConfig.App.LogLevel]
	opts := slog.HandlerOptions{
		Level: logLevel,
	}
	handler := slog.NewTextHandler(os.Stdout, &opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Routes
	handleTempl("/", func() templ.Component {
		forecast, err := yr.NewYRLocationForecast(appConfig.Weather.Lat, appConfig.Weather.Lon)
		if err != nil {
			slog.Error("An error occurred while fetching weather forcasts from the YR.no API",
				"err", err)
			// TODO: UI representation of errors
		}

		model := index.Model{
			Config:     *appConfig,
			Time:       clock.NewModel(),
			YRNow:      yr.NewYRNowModel(*appConfig, *forecast),
			YRForecast: yr.NewYRForecastModel(*appConfig, *forecast),
		}

		return index.View(model, getPageModel("/"))
	})

	handleTempl("/htmx/time", func() templ.Component {
		return clock.View(clock.NewModel())
	})

	handleTempl("/htmx/yrnow", func() templ.Component {
		forecast, err := yr.NewYRLocationForecast(appConfig.Weather.Lat, appConfig.Weather.Lon)
		if err != nil {
			slog.Error("An error occurred while fetching weather forcasts from the YR.no API",
				"err", err)
			// TODO: UI representation of errors
		}
		return yr.YRNowView(yr.NewYRNowModel(*appConfig, *forecast))
	})

	handleTempl("/htmx/yrforecast", func() templ.Component {
		forecast, err := yr.NewYRLocationForecast(appConfig.Weather.Lat, appConfig.Weather.Lon)
		if err != nil {
			slog.Error("An error occurred while fetching weather forcasts from the YR.no API",
				"err", err)
			// TODO: UI representation of errors
		}
		return yr.YRForecastView(yr.NewYRForecastModel(*appConfig, *forecast))
	})

	// Static files
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.Handle("/static/", fs)

	// Run
	port := ":" + strconv.Itoa(appConfig.App.Port)
	slog.Info("Farstu started", "port", port)
	http.ListenAndServe(port, nil)
}
