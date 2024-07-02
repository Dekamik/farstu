package main

import (
	"context"
	"farstu/internal/templates"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func helloHandler(c echo.Context) error {
	return templ.Handler(
		templates.Hello("World")).Component.Render(context.Background(),
		c.Response().Writer,
	)
}

func main() {
	e := echo.New()

	e.GET("/", helloHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
