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

	// Trying to make a slice of an empty array makes the program panic
	selectedDepartures := make([]Departure, 0)
	if len(departures) != 0 {
		selectedDepartures = departures[0:config.SL.MaxRows]
	}

	return DeparturesViewModel{
		Departures: selectedDepartures,
		Enabled:    config.SL.Enabled,
	}
}
