package sl

import (
	"sort"
	"strings"

	"github.com/Dekamik/farstu/internal/config"
)

type Deviation struct {
	Render          DeviationRender
	Priority        DeviationPriority
	MessageVariants map[string]DeviationMessage
	Lines           []DeviationLine
}

type DeviationRender struct {
	Color string
	Modes []DeviationRenderMode
}

type DeviationRenderMode struct {
	Color string
	Mode  string
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

	strToMode := map[string]DeviationRenderMode{
		"buss": {
			Color: "",
			Mode:  "BUS",
		},
		"Pendeltåg": {
			Color: "",
			Mode:  "TRAIN",
		},
		"tunnelbanans gröna linje": {
			Color: "success",
			Mode:  "METRO",
		},
		"tunnelbanans röda linje": {
			Color: "danger",
			Mode:  "METRO",
		},
		"tunnelbanans blåa linje": {
			Color: "primary",
			Mode:  "METRO",
		},
	}

	linesSet := make(map[string]bool)
	for _, line := range deviation.Lines {
		linesSet[line.GroupOfLines] = true
	}

	if strings.Contains(strings.ToLower(deviation.MessageVariants["sv"].ScopeAlias), "buss") {
		linesSet["buss"] = true
	}

	modes := make([]DeviationRenderMode, 0)
	for k := range linesSet {
		if val, ok := strToMode[k]; ok {
			modes = append(modes, val)
		}
	}

	return DeviationRender{
		Color: color,
		Modes: modes,
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
		return deviations[i].Priority.UrgencyLevel > deviations[j].Priority.UrgencyLevel
	})
	sort.Slice(deviations, func(i, j int) bool {
		return deviations[i].Priority.InfluenceLevel > deviations[j].Priority.InfluenceLevel
	})
	sort.Slice(deviations, func(i, j int) bool {
		return deviations[i].Priority.ImportanceLevel > deviations[j].Priority.ImportanceLevel
	})

	return DeviationsViewModel{
		Deviations: deviations,
	}
}
