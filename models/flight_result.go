package models

import "github.com/MaiaraB/travel-plan/ssmodels"

type FlightsResult struct {
	Currency    ssmodels.Currency
	AgentsInfo  []Agent
	InboundLeg  Leg
	OutboundLeg Leg
}
