package sl

import (
	"sort"

	"github.com/Dekamik/farstu/internal/components/shared"
	"github.com/Dekamik/farstu/internal/config"
)

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
	Render          DeviationRender
	Priority        DeviationPriority
	MessageVariants map[string]DeviationMessage
	Lines           []DeviationLine
}

type DeviationRender struct {
	Color string
}

type DeviationPriority struct {
	ImportanceLevel int
	InfluenceLevel  int
	UrgencyLevel    int
}

type DeviationMessage struct {
	Header     string
	Details    string
	ScopeAlias string
	Weblink    string
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
	Page       shared.PageViewModel
}

func calculateRender(deviation Deviation) DeviationRender {
	var color string = ""

	switch {
	case deviation.Priority.ImportanceLevel >= 9:
		color = "border border-3 border-danger-subtle"
	case deviation.Priority.ImportanceLevel >= 7:
		color = "border border-3 border-warning"
	default:
		color = ""
	}

	return DeviationRender{
		Color: color,
	}
}

func NewDeviationsViewModel(config config.AppConfig, response []slDeviationResponse) DeviationsViewModel {
	deviations := make([]Deviation, 0)

	for i, item := range response {
		deviation := Deviation{
			Priority: DeviationPriority{
				ImportanceLevel: item.Priority.ImportanceLevel,
				InfluenceLevel:  item.Priority.InfluenceLevel,
				UrgencyLevel:    item.Priority.UrgencyLevel,
			},
			MessageVariants: make(map[string]DeviationMessage),
			Lines:           make([]DeviationLine, 0),
		}

		for _, message := range response[i].MessageVariants {
			m := DeviationMessage{
				Header:     message.Header,
				Details:    message.Details,
				ScopeAlias: message.ScopeAlias,
				Weblink:    message.Weblink,
			}
			deviation.MessageVariants[message.Language] = m
		}

		for _, line := range response[i].Scope.Lines {
			l := DeviationLine{
				ID:            line.ID,
				Designation:   line.Designation,
				TransportMode: line.TransportMode,
				Name:          line.Name,
				GroupOfLines:  line.GroupOfLines,
			}
			deviation.Lines = append(deviation.Lines, l)
		}

		deviation.Render = calculateRender(deviation)

		deviations = append(deviations, deviation)
	}

	sort.Slice(deviations, func(i, j int) bool {
		return deviations[i].Priority.ImportanceLevel > deviations[j].Priority.ImportanceLevel
	})

	return DeviationsViewModel{
		Deviations: deviations,
	}
}
