package main

import (
	"farstu/internal/views"
	"log"
	"net/http"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/a-h/templ"
)

type Config struct {
	App struct {
		Port int
	}
	GTFS struct {
		Enabled               bool
		RegionalRealtimeKey   string
		RegionalStaticDataKey string
	}
	Weather struct {
		Enabled bool
		Lat     float64
		Lon     float64

		Colors struct {
			TempColorCoolCoolest string
			TempColorCoolHottest string
			TempColorMid         string
			TempColorHotCoolest  string
			TempColorHotHottest  string

			ClassPrecip   string
			ClassNoPrecip string
		}
	}
}

func main() {
	// Config
	var config Config
	_, err := toml.DecodeFile("app.toml", &config)
	if err != nil {
		log.Fatal(err)
	}

	// Routes
	http.Handle("/", templ.Handler(views.Hello("World")))

	// Run
	port := ":" + strconv.Itoa(config.App.Port)
	http.ListenAndServe(port, nil)
}
