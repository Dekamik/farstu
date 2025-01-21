package clock

type ViewModel struct {
	Date string
	Time string
}

func NewViewModel() ViewModel {
	t := GetTime()
	return ViewModel{
		Date: t.Date,
		Time: t.Time,
	}
}
