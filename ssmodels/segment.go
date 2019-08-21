package ssmodels

type Segment struct {
	ID                 int `json:"Id"`
	OriginStation      int
	DestinationStation int
	DepartureDateTime  MyTime
	ArrivalDateTime    MyTime
	Carrier            int
	OperatingCarrier   int
	Duration           int
	FlightNumber       string
	JourneyMode        string
	Directionality     string
}

func SearchSegmentByID(list []Segment, id int) Segment {
	var idElement Segment
	for i := range list {
		currentID := list[i].ID
		if currentID == id {
			idElement = list[i]
			break
		}
	}
	return idElement
}
