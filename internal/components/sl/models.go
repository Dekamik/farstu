package sl

import "github.com/Dekamik/farstu/internal/config"

type Departure struct {
	Destination   string
	DisplayTime   string
	LineNumber    string
	TransportMode string
}

type DeparturesViewModel struct {
	Departures []Departure
	Enabled    bool
	Message    string
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

type Deviation struct {
	ImportanceLevel int
	InfluenceLevel  int
	UrgencyLevel    int
	Header          string
	Details         string
	ScopeAlias      string
	Weblink         string
	Lines           []DeviationLine
}

type DeviationLine struct {
	ID            int
	Designation   string
	TransportMode string
	Name          string
	GroupOfLines  string
}

type DeviationsViewModel struct {
	Deviations []Deviation
	Message    string
}

func NewDeviationsViewModel(config config.AppConfig, response slDeviationsResponse) DeviationsViewModel {
	deviations := make([]Deviation, 0)

	for i, item := range response.Deviations {
		deviation := Deviation{
			ImportanceLevel: item.Priority.ImportanceLevel,
			InfluenceLevel:  item.Priority.InfluenceLevel,
			UrgencyLevel:    item.Priority.UrgencyLevel,
			Header:          item.MessageVariants[0].Header,
			Details:         item.MessageVariants[0].Details,
			ScopeAlias:      item.MessageVariants[0].ScopeAlias,
			Weblink:         item.MessageVariants[0].Weblink,
		}

		for _, line := range response.Deviations[i].Scope.Lines {
			l := DeviationLine{
				ID:            line.ID,
				Designation:   line.Designation,
				TransportMode: line.TransportMode,
				Name:          line.Name,
				GroupOfLines:  line.GroupOfLines,
			}
			deviation.Lines = append(deviation.Lines, l)
		}

		deviations = append(deviations, deviation)
	}

	return DeviationsViewModel{
		Deviations: deviations,
	}
}
