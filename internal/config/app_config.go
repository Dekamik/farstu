package config

import (
	"github.com/BurntSushi/toml"
)

type AppConfig struct {
	App struct {
		Environment string
		LogFile     string
		LogLevel    string
		Port        int
	}
	SL struct {
		Enabled    bool
		MaxRows    int
		SiteName   string
		Deviations struct {
			Future bool
			Lines  []int
			Sites  []int
		}
	}
	Weather struct {
		Enabled bool
		Lat     float64
		Lon     float64
		MaxRows int

		Colors struct {
			TempMin float64
			TempMid float64
			TempMax float64

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

func ReadAppConfig(path string) (*AppConfig, error) {
	var config AppConfig

	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
