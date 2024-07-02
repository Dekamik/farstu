package index

import (
	"farstu/internal/clock"
	"farstu/internal/yr"
)

type Model struct {
	Time clock.Model
	YRNow yr.YRNowModel
}

func NewModel() Model {
	return Model{
		Time: clock.NewModel(),
	}
}
