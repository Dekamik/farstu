package yr

import (
	"farstu/internal/api"
	"fmt"
	"time"
)

type yrLocationForecast struct {
	Type     string `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Meta struct {
			UpdatedAt time.Time `json:"updated_at"`
			Units     struct {
				AirPressureAtSeaLevel      string `json:"air_pressure_at_sea_level"`
				AirTemperature             string `json:"air_temperature"`
				AirTemperatureMax          string `json:"air_temperature_max"`
				AirTemperatureMin          string `json:"air_temperature_min"`
				AirTemperaturePercentile10 string `json:"air_temperature_percentile_10"`
				AirTemperaturePercentile90 string `json:"air_temperature_percentile_90"`
				CloudAreaFraction          string `json:"cloud_area_fraction"`
				CloudAreaFractionHigh      string `json:"cloud_area_fraction_high"`
				CloudAreaFractionLow       string `json:"cloud_area_fraction_low"`
				CloudAreaFractionMedium    string `json:"cloud_area_fraction_medium"`
				DewPointTemperature        string `json:"dew_point_temperature"`
				FogAreaFraction            string `json:"fog_area_fraction"`
				PrecipitationAmount        string `json:"precipitation_amount"`
				PrecipitationAmountMax     string `json:"precipitation_amount_max"`
				PrecipitationAmountMin     string `json:"precipitation_amount_min"`
				ProbabilityOfPrecipitation string `json:"probability_of_precipitation"`
				ProbabilityOfThunder       string `json:"probability_of_thunder"`
				RelativeHumidity           string `json:"relative_humidity"`
				UltravioletIndexClearSky   string `json:"ultraviolet_index_clear_sky"`
				WindFromDirection          string `json:"wind_from_direction"`
				WindSpeed                  string `json:"wind_speed"`
				WindSpeedOfGust            string `json:"wind_speed_of_gust"`
				WindSpeedPercentile10      string `json:"wind_speed_percentile_10"`
				WindSpeedPercentile90      string `json:"wind_speed_percentile_90"`
			} `json:"units"`
		} `json:"meta"`
		Timeseries []struct {
			Time time.Time `json:"time"`
			Data struct {
				Instant struct {
					Details struct {
						AirPressureAtSeaLevel      float64 `json:"air_pressure_at_sea_level"`
						AirTemperature             float64 `json:"air_temperature"`
						AirTemperaturePercentile10 float64 `json:"air_temperature_percentile_10"`
						AirTemperaturePercentile90 float64 `json:"air_temperature_percentile_90"`
						CloudAreaFraction          float64 `json:"cloud_area_fraction"`
						CloudAreaFractionHigh      float64 `json:"cloud_area_fraction_high"`
						CloudAreaFractionLow       float64 `json:"cloud_area_fraction_low"`
						CloudAreaFractionMedium    float64 `json:"cloud_area_fraction_medium"`
						DewPointTemperature        float64 `json:"dew_point_temperature"`
						FogAreaFraction            float64 `json:"fog_area_fraction,omitempty"`
						RelativeHumidity           float64 `json:"relative_humidity"`
						UltravioletIndexClearSky   float64 `json:"ultraviolet_index_clear_sky,omitempty"`
						WindFromDirection          float64 `json:"wind_from_direction"`
						WindSpeed                  float64 `json:"wind_speed"`
						WindSpeedOfGust            float64 `json:"wind_speed_of_gust,omitempty"`
						WindSpeedPercentile10      float64 `json:"wind_speed_percentile_10"`
						WindSpeedPercentile90      float64 `json:"wind_speed_percentile_90"`
					} `json:"details"`
				} `json:"instant"`
				Next12Hours struct {
					Summary struct {
						SymbolCode       string `json:"symbol_code"`
						SymbolConfidence string `json:"symbol_confidence"`
					} `json:"summary"`
					Details struct {
						ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
					} `json:"details"`
				} `json:"next_12_hours,omitempty"`
				Next1Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
					Details struct {
						PrecipitationAmount        float64 `json:"precipitation_amount"`
						PrecipitationAmountMax     float64 `json:"precipitation_amount_max"`
						PrecipitationAmountMin     float64 `json:"precipitation_amount_min"`
						ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
						ProbabilityOfThunder       float64 `json:"probability_of_thunder"`
					} `json:"details"`
				} `json:"next_1_hours,omitempty"`
				Next6Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
					Details struct {
						AirTemperatureMax          float64 `json:"air_temperature_max"`
						AirTemperatureMin          float64 `json:"air_temperature_min"`
						PrecipitationAmount        float64 `json:"precipitation_amount"`
						PrecipitationAmountMax     float64 `json:"precipitation_amount_max"`
						PrecipitationAmountMin     float64 `json:"precipitation_amount_min"`
						ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
					} `json:"details"`
				} `json:"next_6_hours,omitempty"`
			} `json:"data"`
		} `json:"timeseries"`
	} `json:"properties"`
}

func newYRLocationForecast(lat float64, lon float64) (*yrLocationForecast, error) {
	url := fmt.Sprintf("https://api.met.no/weatherapi/locationforecast/2.0/compact?lat=%f&lon=%f", lat, lon)
	return api.GET[yrLocationForecast](url)
}
