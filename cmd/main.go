package main

import (
	"context"
	"farstu/internal/templates"
	"log"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
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

func helloHandler(c echo.Context) error {
	return templ.Handler(
		templates.Hello("World")).Component.Render(context.Background(),
		c.Response().Writer,
	)
}

func main() {
	// Config
	var config Config
	_, err := toml.DecodeFile("app.toml", &config)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	// Routes
	e.GET("/", helloHandler)

	// Run
	port := ":" + strconv.Itoa(config.App.Port)
	e.Logger.Fatal(e.Start(port))
}
