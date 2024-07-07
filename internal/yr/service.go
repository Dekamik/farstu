package yr

import "farstu/internal/cache"

type YRService interface {
	GetForecast() (*YRLocationForecast, error)
}

type yrServiceImpl struct {
	cachedForecast cache.Cache[YRLocationForecast, refreshForecastArgs]
}

var _ YRService = yrServiceImpl{}

func (y yrServiceImpl) GetForecast() (*YRLocationForecast, error) {
	return y.cachedForecast.GetOrRefresh()
}

type refreshForecastArgs struct {
	Lat float64
	Lon float64
}

func refreshForecast(args refreshForecastArgs) (*YRLocationForecast, error) {
	return NewYRLocationForecast(args.Lat, args.Lon)
}

type YRServiceArgs struct {
	ForecastTTL int
	Lat         float64
	Lon         float64
}

func NewYRService(args YRServiceArgs) YRService {
	refreshForecastArgs := refreshForecastArgs{
		Lat: args.Lat,
		Lon: args.Lon,
	}
	return &yrServiceImpl{
		cachedForecast: cache.New(args.ForecastTTL, refreshForecast, refreshForecastArgs),
	}
}
