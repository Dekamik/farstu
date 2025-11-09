package index

import (
	"net/http"

	"github.com/Dekamik/farstu/internal/routes/shared"
	"github.com/Dekamik/farstu/internal/routes/index/components/sl"
)

type Index struct {
	Departures []sl.Departure
}

func Handler(w http.ResponseWriter, _ *http.Request) {
	data := Index{
		Departures: []sl.Departure{
			{TransportMode: "BUS", LineNumber: "504", Destination: "Sundbybergs Station", DisplayTime: "2 min"},
			{TransportMode: "BUS", LineNumber: "504", Destination: "Rissne", DisplayTime: "4 min"},
			{TransportMode: "BUS", LineNumber: "504", Destination: "Sundbybergs Station", DisplayTime: "9 min"},
			{TransportMode: "BUS", LineNumber: "504", Destination: "Rissne", DisplayTime: "11 min"},
			{TransportMode: "BUS", LineNumber: "504", Destination: "Sundbybergs Station", DisplayTime: "16 min"},
			{TransportMode: "BUS", LineNumber: "504", Destination: "Rissne", DisplayTime: "18 min"},
		},
	}
	shared.ExecuteLayout(w, "/", data, "internal/components/index/index.html")
}
