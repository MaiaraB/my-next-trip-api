package ssmodels

type PricingOption struct {
	AgentIds          []int `json:"Agents"`
	QuoteAgeInMinutes int
	Price             float64
	DeeplinkURL       string `json:"DeepLinkUrl"`
}
