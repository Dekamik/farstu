package yr

import (
	"fmt"
	"github.com/Dekamik/farstu/internal/asserts"
	"github.com/Dekamik/farstu/internal/config"
	"image/color"
	"log/slog"
	"math"
)

func ParseHexToColor(s string) (c color.RGBA, err error) {
	asserts.Assert(s[0] == '#', "colors must begin with #")
	asserts.Assert(len(s) == 7 || len(s) == 4, "colors must be formatted like #FFFFFF or #FFF")

	c.A = 0xff

	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)

	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits
		c.R *= 17
		c.G *= 17
		c.B *= 17

	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")
	}

	return
}

func HexString(c color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

func Lerp(min color.RGBA, max color.RGBA, value float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(min.R)*(1-value) + float64(max.R)*value),
		G: uint8(float64(min.G)*(1-value) + float64(max.G)*value),
		B: uint8(float64(min.B)*(1-value) + float64(max.B)*value),
		A: 0xff,
	}
}

func LerpHexString(min string, max string, value float64) (string, error) {
	minC, err := ParseHexToColor(min)
	if err != nil {
		return "", err
	}

	maxC, err := ParseHexToColor(max)
	if err != nil {
		return "", err
	}

	c := Lerp(minC, maxC, value)

	return HexString(c), nil
}

func getTemperatureColor(conf config.AppConfig, temperature float64) string {
	tempMin := conf.Weather.Colors.TempMin
	tempMid := conf.Weather.Colors.TempMid
	tempMax := conf.Weather.Colors.TempMax

	asserts.Assert(tempMin < tempMid, "Min must be smaller than mid")
	asserts.Assert(tempMid < tempMax, "Mid must be smaller than max")
	asserts.Assert(tempMin < tempMax, "Min must be smaller than max")

	if temperature == tempMid {
		return conf.Weather.Colors.TempColorMid
	} else if temperature <= tempMin {
		return conf.Weather.Colors.TempColorCoolCoolest
	} else if temperature >= tempMax {
		return conf.Weather.Colors.TempColorHotHottest
	}

	var color0 string
	var color1 string
	var value float64

	if temperature > tempMid {
		color0 = conf.Weather.Colors.TempColorHotCoolest
		color1 = conf.Weather.Colors.TempColorHotHottest
		value = math.Abs(temperature / tempMax)
	} else {
		color0 = conf.Weather.Colors.TempColorCoolHottest
		color1 = conf.Weather.Colors.TempColorCoolCoolest
		value = math.Abs(temperature / tempMin)
	}

	color, err := LerpHexString(color0, color1, value)
	if err != nil {
		slog.Error("error occurred when lerping colors", "err", err)
		return conf.Weather.Colors.TempColorMid
	}

	return color
}

func getPrecipitationColorClass(conf config.AppConfig, precipitation float64) string {
	if precipitation > 0 {
		return conf.Weather.Colors.ClassPrecip
	}
	return conf.Weather.Colors.ClassNoPrecip
}
