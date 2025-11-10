package index

import (
	"github.com/Dekamik/farstu/internal/routes/index/components/sl"
	"github.com/Dekamik/farstu/internal/routes/index/components/yr"
)

type Index struct {
	Departures []sl.Departure
	Forecast   []yr.YRForecastItem
}
