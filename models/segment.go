package models

import "github.com/MaiaraB/travel-plan/ssmodels"

type Segment struct {
	Origin      Place
	Destination Place
	Departure   ssmodels.MyTime
	Arrival     ssmodels.MyTime
	Duration    int
}
