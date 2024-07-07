package main

import (
	"farstu/internal/clock"
	"farstu/internal/config"
	"farstu/internal/gtfs"
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

var yrService yr.YRService

func handleTempl(path string, componentFunc func() templ.Component) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		componentFunc().Render(r.Context(), w)
	})
}

func getWeatherViewModels(appConfig config.AppConfig) (yr.YRNowViewModel, yr.YRForecastViewModel) {
	var yrNowViewModel yr.YRNowViewModel
	var yrForecastViewModel yr.YRForecastViewModel
	
	forecast, err := yrService.GetForecast()
	if err != nil {
		slog.Error("An error occurred while fetching weather forcasts from the YR.no API",
			"err", err)
		yrNowViewModel = yr.YRNowViewModel{
			Enabled: appConfig.Weather.Enabled,
			Message: "Fel vid väderdatahämtning",
		}
		yrForecastViewModel = yr.YRForecastViewModel{
			Enabled: appConfig.Weather.Enabled,
			Message: "Fel vid väderdatahämtning",
		}
	} else {
		yrNowViewModel = yr.NewYRNowViewModel(appConfig, *forecast)
		yrForecastViewModel = yr.NewYRForecastViewModel(appConfig, *forecast)
	}

	return yrNowViewModel, yrForecastViewModel
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

	// Services
	yrServiceArgs := yr.YRServiceArgs{
		ForecastTTL: 300,
		Lat: appConfig.Weather.Lat,
		Lon: appConfig.Weather.Lon,
	}
	yrService = yr.NewYRService(yrServiceArgs)

	// Routes
	handleTempl("/", func() templ.Component {
		yrNowViewModel, yrForecastViewModel := getWeatherViewModels(*appConfig)

		model := index.ViewModel{
			Config:     *appConfig,
			Departures: gtfs.NewDeparturesViewModel(),
			Time:       clock.NewViewModel(),
			YRNow:      yrNowViewModel,
			YRForecast: yrForecastViewModel,
		}

		return index.View(model, shared.NewPageViewModel("/"))
	})

	handleTempl("/htmx/departures", func() templ.Component {
		return gtfs.DeparturesView(gtfs.NewDeparturesViewModel())
	})

	handleTempl("/htmx/time", func() templ.Component {
		return clock.View(clock.NewViewModel())
	})

	handleTempl("/htmx/yrnow", func() templ.Component {
		yrNowViewModel, _ := getWeatherViewModels(*appConfig)
		return yr.YRNowView(yrNowViewModel)
	})

	handleTempl("/htmx/yrforecast", func() templ.Component {
		_, yrForecastViewModel := getWeatherViewModels(*appConfig)
		return yr.YRForecastView(yrForecastViewModel)
	})

	// Static files
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.Handle("/static/", fs)

	// Run
	port := ":" + strconv.Itoa(appConfig.App.Port)
	slog.Info("Farstu started", "port", port)
	http.ListenAndServe(port, nil)
}
