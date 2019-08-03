package main

import (
	"log"
	"time"
)

type Interval struct {
	Outbound time.Time
	Inbound  time.Time
}

func getDateIntervals(weekday time.Weekday, duration int, beginDate time.Time, endDate time.Time) []Interval {
	var intervals []Interval

	outbound := beginDate

	for i := 1; outbound.Weekday() != weekday; i++ {
		outbound = beginDate.AddDate(0, 0, i)
	}

	inbound := outbound.AddDate(0, 0, duration)

	for {
		if inbound.Before(endDate) {
			intervals = append(intervals, Interval{Outbound: outbound, Inbound: inbound})
		} else {
			break
		}
		outbound = outbound.AddDate(0, 0, 7)
		inbound = outbound.AddDate(0, 0, duration)
		log.Println("OUTBOUND: ", outbound, " / INBOUND: ", inbound)
	}

	return intervals
}
