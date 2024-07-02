# Farstu

Farstu is a digital signage application next to your door. It displays a weather
forecast and upcoming public transport departures.

# Getting started

## Building the application

1. Install Go by following the instructions for your OS at
   [go.dev](https://go.dev/doc/install)
2. Run `make deps` to install dev dependencies
3. Run `make` to generate templates, install modules and build binary
4. Run `make install` to install the app on your system

## Configuring the application

Copy `example.app.toml` to `app.toml` and configure GTFS keys and weather
lat/lon inside `app.toml`.

You can get the GTFS keys at [Trafiklab](https://www.trafiklab.se/api).

## Running the application

Stand in the directory where you put your `app.toml` and run `farstu`. The app
should now be running at the configured port (default `http://localhost:8080`).

### Development

You can run the application directly by running `go run cmd/main.go`.
