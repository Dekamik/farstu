package clock

type Model struct {
	Date string
	Time string
}

func NewModel() Model {
	t := GetTime()
	return Model{
		Date: t.Date,
		Time: t.Time,
	}
}
