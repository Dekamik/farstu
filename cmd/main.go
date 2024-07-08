package main

import (
	"farstu/internal/clock"
	"farstu/internal/config"
	"farstu/internal/index"
	"farstu/internal/shared"
	"farstu/internal/sl"
	"farstu/internal/yr"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

var logLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

var slService sl.SLService
var yrService yr.YRService

func handleTempl(path string, componentFunc func() templ.Component) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		componentFunc().Render(r.Context(), w)
	})
}

func getDeparturesViewModel(appConfig config.AppConfig) sl.DeparturesViewModel {
	var slDeparturesViewModel sl.DeparturesViewModel

	departures, err := slService.GetDepartures()
	if err != nil {
		slog.Warn("an error occurred when fetching departures from SL", "err", err)
		slDeparturesViewModel = sl.DeparturesViewModel{
			Enabled: appConfig.SL.Enabled,
			Message: "Fel vid avgångsdatahämtning",
		}
	} else {
		slDeparturesViewModel = sl.NewDeparturesViewModel(appConfig, *departures)
	}

	return slDeparturesViewModel
}

func getWeatherViewModels(appConfig config.AppConfig) (yr.YRNowViewModel, yr.YRForecastViewModel) {
	var yrNowViewModel yr.YRNowViewModel
	var yrForecastViewModel yr.YRForecastViewModel

	forecast, err := yrService.GetForecast()
	if err != nil {
		slog.Warn("an error occurred when fetching weather forcasts from YR.no", "err", err)
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
	var err error

	// Config
	appConfigPath := "app.toml"
	appConfig, err := config.ReadAppConfig(appConfigPath)
	if err != nil {
		slog.Error("An error occurred while reading "+appConfigPath, "err", err)
		os.Exit(1)
	}

	// Logging
	levelSanitized := strings.ToLower(appConfig.App.LogLevel)
	logLevel := logLevelMap[levelSanitized]
	opts := slog.HandlerOptions{ Level: logLevel }
	writers := []io.Writer { os.Stdout }

	if appConfig.App.LogFile != "" {
		errWhenCreatingDirectory := false
		path := filepath.Dir(appConfig.App.LogFile)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				slog.Warn("an error occurred when creating log directory")
				errWhenCreatingDirectory = true
			}
		}

		if !errWhenCreatingDirectory {
			fWriter, err := os.OpenFile(appConfig.App.LogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				slog.Warn("an error occurred when opening log file", "err", err)
			} else {
				writers = append(writers, fWriter)
			}
		}
	}

	multiWriter := io.MultiWriter(writers...)
	handler := slog.NewTextHandler(multiWriter, &opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Services
	slServiceArgs := sl.SLServiceArgs{
		DeparturesTTL: 300,
		SiteName:      appConfig.SL.SiteName,
	}
	slService, err = sl.NewSLService(slServiceArgs)
	if err != nil {
		if err == sl.ErrSiteIDNotFound {
			slog.Error("SL site ID not found", "site_name", appConfig.SL.SiteName)
			os.Exit(1)
		} else {
			slog.Error("an unhandled error occurred", "err", err)
			os.Exit(1)
		}
	}

	yrServiceArgs := yr.YRServiceArgs{
		ForecastTTL: 300,
		Lat:         appConfig.Weather.Lat,
		Lon:         appConfig.Weather.Lon,
	}
	yrService = yr.NewYRService(yrServiceArgs)

	// Routes
	handleTempl("/", func() templ.Component {
		yrNowViewModel, yrForecastViewModel := getWeatherViewModels(*appConfig)

		model := index.ViewModel{
			Config:     *appConfig,
			Departures: getDeparturesViewModel(*appConfig),
			Time:       clock.NewViewModel(),
			YRNow:      yrNowViewModel,
			YRForecast: yrForecastViewModel,
		}

		return index.View(model, shared.NewPageViewModel("/"))
	})

	handleTempl("/htmx/departures", func() templ.Component {
		return sl.DeparturesView(getDeparturesViewModel(*appConfig))
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
