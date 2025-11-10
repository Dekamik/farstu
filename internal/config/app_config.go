package config

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	App struct {
		Environment string `json:"environment"`
		LogFile     string `json:"logFile"`
		LogLevel    string `json:"logLevel"`
		Port        int    `json:"port"`
	} `json:"app"`
	SL struct {
		SiteName   string `json:"siteName"`
		Deviations struct {
			Future bool  `json:"future"`
			Lines  []int `json:"lines"`
			Sites  []int `json:"sites"`
		} `json:"deviations"`
	} `json:"sl"`
	Weather struct {
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
		Colors struct {
			TempMin              float64 `json:"tempMin"`
			TempMid              float64 `json:"tempMid"`
			TempMax              float64 `json:"tempMax"`
			TempColorCoolCoolest string  `json:"tempColorCoolCoolest"`
			TempColorCoolHottest string  `json:"tempColorCoolHottest"`
			TempColorMid         string  `json:"tempColorMid"`
			TempColorHotCoolest  string  `json:"tempColorHotCoolest"`
			TempColorHotHottest  string  `json:"tempColorHotHottest"`
			ClassPrecip          string  `json:"classPrecip"`
			ClassNoPrecip        string  `json:"classNoPrecip"`
		} `json:"colors"`
	} `json:"weather"`
}

func ReadAppConfig(path string) (*AppConfig, error) {
	var config AppConfig

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
