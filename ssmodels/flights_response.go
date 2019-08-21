package ssmodels

type FlightsResponse struct {
	Query       Query
	Status      string
	Itineraries []Itinerary
	Legs        []Leg
	Segments    []Segment
	Carriers    []Carrier
	Agents      []Agent
	Places      []Place
	Currencies  []Currency
}
