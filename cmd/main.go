package main

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Dekamik/farstu/internal/config"
	"github.com/Dekamik/farstu/internal/routes/deviations"
	"github.com/Dekamik/farstu/internal/routes/index"
	"github.com/Dekamik/farstu/internal/routes/index/components/sl"
	"github.com/Dekamik/farstu/internal/routes/index/components/yr"
	"github.com/Dekamik/farstu/internal/routes/settings"
	"github.com/Dekamik/farstu/internal/routes/shared"
)

var logLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func main() {
	// Config
	appConfigPath := os.Getenv("FARSTU_APP_CONFIG_PATH")
	if appConfigPath == "" {
		appConfigPath = "app.json"
	}

	appConfig, err := config.Read(appConfigPath)
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
				slog.Warn("an error occurred when creating log directory", "err", err)
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
	var slService sl.SLService
	var yrService yr.YRService

	slServiceArgs := sl.SLServiceArgs{
		DeparturesTTL:  15,
		DeviationsTTL:  3,
		InitRetriesSec: []int{1, 4, 8, 8},
	}
	slService, err = sl.NewSLService(slServiceArgs, appConfigPath)
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
	yrService, err = yr.NewYRService(yrServiceArgs, appConfigPath)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c, err := config.Read(appConfigPath)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		data := index.Index{
			Departures: slService.GetDepartures(),
			Forecast: yrService.GetForecast(*appConfig),
		}
		templates := []string{
			"internal/routes/index/index.html",
			"internal/routes/index/components/sl/sl.html",
			"internal/routes/index/components/yr/yr.html",
		}
		shared.ExecuteLayout(w, "/", c, data, templates...)
	})

	http.HandleFunc("/deviations", func(w http.ResponseWriter, r *http.Request) {
		c, err := config.Read(appConfigPath)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		data := deviations.Deviations{
			Deviations: slService.GetDeviations(),
		}
		shared.ExecuteLayout(w, "deviations", c, data, "internal/routes/deviations/deviations.html")
	})

	http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		c, err := config.Read(appConfigPath)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		data := settings.Settings{}
		shared.ExecuteLayout(w, "settings", c, data, "internal/routes/settings/settings.html")
	})

	http.HandleFunc("/htmx/time", func(w http.ResponseWriter, r *http.Request) {
		data := shared.GetTime()
		tmpl := template.Must(template.ParseFiles("internal/routes/shared/layout/clock.html"))
		err := tmpl.Execute(w, data)
		if err != nil {
			slog.Error(err.Error())
		}
	})

	http.HandleFunc("/htmx/sl", func(w http.ResponseWriter, r *http.Request) {
		data := index.Index{
			Departures: slService.GetDepartures(),
		}
		tmpl := template.Must(template.ParseFiles("internal/routes/index/components/sl/sl.html"))
		err := tmpl.Execute(w, data)
		if err != nil {
			slog.Error(err.Error())
		}
	})

	http.HandleFunc("/htmx/yr", func(w http.ResponseWriter, r *http.Request) {
		c, err := config.Read(appConfigPath)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		data := index.Index{
			Forecast: yrService.GetForecast(*c),
		}
		tmpl := template.Must(template.ParseFiles("internal/routes/index/components/yr/yr.html"))
		err = tmpl.Execute(w, data)
		if err != nil {
			slog.Error(err.Error())
		}
	})

	// Static files
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.Handle("/static/", fs)

	// Run
	port := ":" + strconv.Itoa(appConfig.App.Port)
	slog.Info("Farstu started", "port", port)
	http.ListenAndServe(port, nil)
}
