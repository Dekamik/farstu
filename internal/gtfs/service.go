package gtfs

type GTFSService interface{}

type gtfsServiceImpl struct{}

var _ GTFSService = gtfsServiceImpl{}

func NewGTFSService() GTFSService {
	return &gtfsServiceImpl{}
}
