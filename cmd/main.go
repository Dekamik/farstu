package main

import (
	"github.com/Dekamik/farstu/internal/components/clock"
	"github.com/Dekamik/farstu/internal/components/index"
	"github.com/Dekamik/farstu/internal/components/page"
	"github.com/Dekamik/farstu/internal/components/sl"
	"github.com/Dekamik/farstu/internal/components/yr"
	"github.com/Dekamik/farstu/internal/config"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
)

var logLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func handleTempl(path string, componentFunc func() templ.Component) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		componentFunc().Render(r.Context(), w)
	})
}

func main() {
	var slService sl.SLService
	var yrService yr.YRService
	var err error

	// Config
	appConfigPath := "app.toml"
	appConfig, err := config.ReadAppConfig(appConfigPath)
	if err != nil {
		slog.Error("An error occurred while reading "+appConfigPath, "err", err)
		os.Exit(1)
	}

	err = os.Setenv("FARSTU_ENVIRONMENT", appConfig.App.Environment)
	if err != nil {
		slog.Error("An error occurred when setting FARSTU_ENVIRONMENT="+appConfig.App.Environment, "err", err)
		os.Exit(1)
	}

	// Logging
	levelSanitized := strings.ToLower(appConfig.App.LogLevel)
	logLevel := logLevelMap[levelSanitized]
	opts := slog.HandlerOptions{Level: logLevel}
	writers := []io.Writer{os.Stdout}

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
		DeparturesTTL:  300,
		InitRetriesSec: []int{1, 4, 8, 8},
		SiteName:       appConfig.SL.SiteName,
	}
	slService, err = sl.NewSLService(slServiceArgs, *appConfig)
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
	yrService = yr.NewYRService(yrServiceArgs, *appConfig)

	// Routes
	handleTempl("/", func() templ.Component {
		yrNowViewModel, yrForecastViewModel := yrService.GetViewModels()

		model := index.ViewModel{
			Config:     *appConfig,
			Departures: slService.GetViewModel(),
			Time:       clock.NewViewModel(),
			YRNow:      yrNowViewModel,
			YRForecast: yrForecastViewModel,
		}

		seasonAndTimeOfDay, err := page.GetSeasonAndTimeOfDay(*appConfig)
		if err != nil {
			slog.Error("an unhandled error occurred", "err", err)
			// Try again in an hour
			nextRetry := int(time.Now().Local().Add(time.Duration(1) * time.Hour).Sub(time.Now().Local()).Seconds())
			seasonAndTimeOfDay = &page.SeasonAndTimeOfDay{
				Season:                   "summer",
				SecondsUntilNextSunEvent: nextRetry,
				TimeOfDay:                "day",
			}
		}

		args := page.NewPageViewModelArgs{
			ActiveHref:               "/",
			SecondsUntilNextSunEvent: seasonAndTimeOfDay.SecondsUntilNextSunEvent,
			Season:                   seasonAndTimeOfDay.Season,
			TimeOfDay:                seasonAndTimeOfDay.TimeOfDay,
		}
		return index.View(model, page.NewPageViewModel(args))
	})

	handleTempl("/htmx/departures", func() templ.Component {
		return sl.DeparturesView(slService.GetViewModel())
	})

	handleTempl("/htmx/departures/next", func() templ.Component {
		return sl.SLDeparturesNextView(slService.GetViewModel())
	})

	handleTempl("/htmx/departures/list", func() templ.Component {
		return sl.SLDeparturesListView(slService.GetViewModel())
	})

	handleTempl("/htmx/time", func() templ.Component {
		return clock.View(clock.NewViewModel())
	})

	handleTempl("/htmx/yrnow", func() templ.Component {
		yrNowViewModel, _ := yrService.GetViewModels()
		return yr.YRNowView(yrNowViewModel)
	})

	handleTempl("/htmx/yrforecast", func() templ.Component {
		_, yrForecastViewModel := yrService.GetViewModels()
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
