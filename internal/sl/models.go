package sl

import "farstu/internal/config"

type DeparturesViewModel struct {
	Departures []Departure
	Enabled    bool
	Message    string
}

type Departure struct {
	Destination   string
	DisplayTime   string
	LineNumber    string
	TransportMode string
}

func NewDeparturesViewModel(config config.AppConfig, response slSiteDeparturesResponse) DeparturesViewModel {
	departures := make([]Departure, 0)

	for _, item := range response.Departures {
		departure := Departure{
			Destination:   item.Destination,
			DisplayTime:   item.Display,
			LineNumber:    item.Line.Designation,
			TransportMode: item.Line.TransportMode,
		}
		departures = append(departures, departure)
	}

	// Bounds checking required, or it will crash
	selectedDepartures := make([]Departure, 0)
	if len(departures) > 0 && len(departures) > config.SL.MaxRows {
		selectedDepartures = departures[0:config.SL.MaxRows]
	} else {
		selectedDepartures = departures
	}

	return DeparturesViewModel{
		Departures: selectedDepartures,
		Enabled:    config.SL.Enabled,
	}
}
