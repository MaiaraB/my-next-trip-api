package ssmodels

type Leg struct {
	ID                  string `json:"Id"`
	SegmentIds          []int
	OriginStation       int
	DestinationStation  int
	Departure           MyTime
	Arrival             MyTime
	Duration            int
	JourneyMode         string
	StopsIds            []int `json:"Stops"`
	CarrierIds          []int `json:"Carriers"`
	OperatingCarrierIds []int `json:"OperatingCarriers"`
	Directionality      string
	FlightNumbers       []FlightNumber
}

func SearchLegByID(list []Leg, id string) Leg {
	var idElement Leg
	for i := range list {
		currentID := list[i].ID
		if currentID == id {
			idElement = list[i]
			break
		}
	}
	return idElement
}
