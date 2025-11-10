package sl

import (
	"strings"
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
