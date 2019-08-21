package models

import "github.com/MaiaraB/travel-plan/ssmodels"

type Leg struct {
	Departure   ssmodels.MyTime
	Arrival     ssmodels.MyTime
	Duration    int
	Stops       []Place // place codes
	Origin      Place   // place code
	Destination Place   //place code
	Carriers    []Carrier
	Segments    []Segment
}
