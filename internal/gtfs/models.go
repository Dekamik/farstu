package gtfs

type DeparturesViewModel struct {
	Departures []Departure
	Message    string
}

type Departure struct {
	Destination   string
	DisplayTime   string
	Line          string
	TransportMode string
}

func NewDeparturesViewModel() DeparturesViewModel {
	return DeparturesViewModel{
		Departures: make([]Departure, 0),
		Message: "No data",
	}
}
