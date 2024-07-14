# Farstu

Farstu is a digital signage application next to your door. It displays a weather
forecast and upcoming public transport departures.

# Getting started

## Building the application

1. Install Go by following the instructions for your OS at
   [go.dev](https://go.dev/doc/install)
2. Run `make deps` to install dev dependencies 
3. Run `make` to generate templates, install modules and build binary (note: if
   you cannot run the dependencies (e.g. templ), you're probably missing
`~/go/bin` in your PATH variable)
4. Run `make install` to install the app on your system

## Configuring the application

1. Copy `example.app.toml` to `app.toml`
2. To configure the weather forecasts, go to [latlong.net](https://latlong.net)
   and lookup your address
3. Paste the lat lon coordinates to the `Lat` and `Lon` variables in `app.toml`

If you've followed the steps above, you should be able to run Farstu.

## Running the application

Stand in the directory where the `app.toml` and the `static/` directory is
located and run `farstu`. The app should now be running at the configured port
(default `http://localhost:8080`).

## Running the application on kiosk mode on Raspberry Pi

1. Install unclutter with `sudo apt install unclutter`
2. Copy the contents of the autostart file into either: a. The local autostart
   (`/home/pi/.config/lxsession/LXDE-pi/autostart`) b. The global autostart
(`/etc/xdg/lxsession/LXDE-pi/autostart`)
3. Restart the raspberry pi

### Development

You can run the application directly by running `go run cmd/main.go`.
