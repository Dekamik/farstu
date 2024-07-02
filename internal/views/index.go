package views

type IndexModel struct {
	Time TimeModel
	YRNow YRNowModel
}

func GetIndexModel() IndexModel {
	return IndexModel{
		Time: GetTimeModel(),
	}
}
