package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Dekamik/farstu/internal/components/sl"
	"github.com/Dekamik/farstu/internal/components/yr"
	"github.com/Dekamik/farstu/internal/config"
	"github.com/Dekamik/farstu/internal/routing"
)

var logLevelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func main() {
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
	var slService sl.SLService
	var yrService yr.YRService

	slServiceArgs := sl.SLServiceArgs{
		DeparturesTTL:  15,
		DeviationsTTL:  3,
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

	services := routing.Services{
		AppConfig: appConfig,
		SL: &slService,
		YR: &yrService,
	}
	routing.Routes(services)
	routing.HTMXRoutes(services)

	// Static files
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.Handle("/static/", fs)

	// Run
	port := ":" + strconv.Itoa(appConfig.App.Port)
	slog.Info("Farstu started", "port", port)
	http.ListenAndServe(port, nil)
}
