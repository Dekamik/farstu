# Farstu

Farstu is a digital signage application next to your door. It displays a weather
forecast and upcoming public transport departures.

# Getting started

## Building application

### Install go

Follow the instructions for your OS at [go.dev](https://go.dev/doc/install).

### Install dev dependencies

Run `make deps`

### Generate templ, install modules and build binary

Run `make`

### Install application on your system

Run `make install`

## Configuring application

Copy `example.app.toml` to `app.toml` and configure GTFS keys and weather
lat/lon.

You can get the GTFS keys at [Trafiklab](https://www.trafiklab.se/api).

## Running the application

### Production

Stand in the directory where you put your `app.toml` and run `farstu`. The app
should now be running at the configured port (default `http://localhost:8080`).

### Development

You can run the application directly by running `go run cmd/main.go`
