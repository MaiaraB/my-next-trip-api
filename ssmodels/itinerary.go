package ssmodels

type Itinerary struct {
	OutboundLegID  string `json:"OutboundLegId"`
	InboundLegID   string `json:"InboundLegId"`
	PricingOptions []PricingOption
}
