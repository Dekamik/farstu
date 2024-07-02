package models

type IndexModel struct {
	Time TimeModel
	YRNow YRNowModel
}

func NewIndexModel() IndexModel {
	return IndexModel{
		Time: NewTimeModel(),
	}
}
