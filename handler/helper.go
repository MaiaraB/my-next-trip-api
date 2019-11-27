package handler

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/MaiaraB/travel-plan/models"
	"github.com/MaiaraB/travel-plan/ssmodels"
	"github.com/matryer/runner"
	"github.com/pkg/errors"
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
		if outbound.Before(endDate) || outbound.Equal(endDate) {
			intervals = append(intervals, Interval{Outbound: outbound, Inbound: inbound})
			log.Println("OUTBOUND: ", outbound, " / INBOUND: ", inbound)
		} else {
			break
		}
		outbound = outbound.AddDate(0, 0, 7)
		inbound = outbound.AddDate(0, 0, duration)
	}

	return intervals
}

func createDataAndIntervalsForSkyscannerAPI(r *http.Request) (url.Values, []Interval, error) {
	queryValues := r.URL.Query()
	cabinClass := queryValues.Get("cabinClass")
	origin := queryValues.Get("origin")
	destination := queryValues.Get("destination")
	outboundWeekDay, _ := strconv.Atoi(queryValues.Get("outboundWeekDay"))
	duration, _ := strconv.Atoi(queryValues.Get("duration"))
	adults, _ := strconv.Atoi(queryValues.Get("adults"))
	children, _ := strconv.Atoi(queryValues.Get("children"))
	infants, _ := strconv.Atoi(queryValues.Get("infants"))
	if queryValues.Get("duration") == "" {
		log.Println("ONE WAY TRIP")
		duration = -1
	}
	country := queryValues.Get("country")
	currency := queryValues.Get("currency")
	locale := queryValues.Get("locale")
	fromDate := queryValues.Get("fromDate")
	toDate := queryValues.Get("toDate")

	parsedFromDate, err := time.Parse(layoutISO, fromDate)
	if err != nil {
		log.Printf("Error while parsing date: %s", err)
	}
	parsedToDate, err := time.Parse(layoutISO, toDate)
	if err != nil {
		log.Printf("Error while parsing date: %s", err)
	}

	// Checking if the search interval exceeds 3 months
	if parsedToDate.After(parsedFromDate.AddDate(0, 3, 0)) {
		return nil, nil, errors.New("SearchIntervalTooBig")
	}

	log.Println("PARSED FROM DATE: ", parsedFromDate, ", PARSED TO DATE: ", parsedToDate)
	intervals := getDateIntervals(time.Weekday(outboundWeekDay), duration-1, parsedFromDate, parsedToDate)
	log.Println("INTERVALS: ", intervals)

	data := url.Values{}
	data.Set("cabinClass", cabinClass)
	data.Set("adults", strconv.Itoa(adults))
	data.Set("locale", locale)
	data.Set("currency", currency)
	data.Set("country", country)
	data.Set("originPlace", origin)
	data.Set("destinationPlace", destination)
	if children > 0 {
		data.Set("children", strconv.Itoa(children))
	}
	if infants > 0 {
		data.Set("infants", strconv.Itoa(infants))
	}

	log.Printf("CABIN CLASS: %s, ADULTS: %s, CHILDREN: %s, INFANTS: %s", cabinClass, strconv.Itoa(adults), strconv.Itoa(children), strconv.Itoa(infants))

	return data, intervals, nil
}

func getPartialResults(ctx context.Context, respond chan<- []models.FlightsResult, data url.Values, index int) {
	defer func() {
		// log.Printf("Partial results #%d complete", index)
	}()

	flightsResponse := make(chan ssmodels.FlightsResponse)

	task := runner.Go(func(shouldStop runner.S) error {
		getFlightsSkyscanner(shouldStop, flightsResponse, data, index)
		return nil
	})

	select {
	case <-ctx.Done():
		log.Printf("getPartialResults #%d: time to return", index)
		task.Stop()
		select {
		case <-task.StopChan():
			// task successfully stopped
		case <-time.After(1 * time.Second):
			// task didn't stop in time
		}

		// execution continues once the code has stopped or has
		// timed out.
		if task.Err() != nil {
			log.Fatalf("getPartialResults %d failed: %s", index, task.Err())
		}
	case flights := <-flightsResponse:
		var flightResults []models.FlightsResult
		log.Printf("Thread %d: %+v. Size: %d", index, flights.Query, len(flights.Itineraries))
		for _, it := range flights.Itineraries {
			result := models.FlightsResult{}

			// Setting result Currency
			result.Currency = ssmodels.SearchCurrencyByCode(flights.Currencies, flights.Query.Currency)

			// Setting result AgentInfo
			var agentsInfo []models.Agent
			for _, po := range it.PricingOptions {
				agent := ssmodels.SearchAgentByID(flights.Agents, po.AgentIds[0])
				agentsInfo = append(agentsInfo, models.Agent{
					Name:        agent.Name,
					ImageURL:    agent.ImageURL,
					Price:       po.Price,
					DeepLinkURL: po.DeeplinkURL})
			}
			result.AgentsInfo = agentsInfo

			// Setting result OutboundLeg
			result.OutboundLeg = configResponseLeg(flights, it.OutboundLegID)

			// Setting result InboundLeg
			if data.Get("inboundDate") != "" {
				result.InboundLeg = configResponseLeg(flights, it.InboundLegID)
			}

			flightResults = append(flightResults, result)

			log.Println("Itinerary price: ", result.AgentsInfo[0].Price)
		}

		respond <- flightResults
	}
}

func copyValues(data url.Values) url.Values {
	newData := url.Values{}
	for k, v := range data {
		newData[k] = v
	}
	return newData
}

func configResponseLeg(flightResponse ssmodels.FlightsResponse, legID string) models.Leg {

	leg := ssmodels.SearchLegByID(flightResponse.Legs, legID)

	stops := []models.Place{}
	for _, stopID := range leg.StopsIds {
		stop := ssmodels.SearchPlaceByID(flightResponse.Places, stopID)
		stops = append(stops, models.Place{Name: stop.Name, Code: stop.Code})
	}

	origin := ssmodels.SearchPlaceByID(flightResponse.Places, leg.OriginStation)
	destination := ssmodels.SearchPlaceByID(flightResponse.Places, leg.DestinationStation)

	carriers := []models.Carrier{}
	for _, carrierID := range leg.CarrierIds {
		carrier := ssmodels.SearchCarrierByID(flightResponse.Carriers, carrierID)
		carriers = append(carriers, models.Carrier{Name: carrier.Name, ImageURL: carrier.ImageURL})
	}

	segs := []models.Segment{}
	for _, segID := range leg.SegmentIds {
		seg := ssmodels.SearchSegmentByID(flightResponse.Segments, segID)
		origin := ssmodels.SearchPlaceByID(flightResponse.Places, seg.OriginStation)
		destination := ssmodels.SearchPlaceByID(flightResponse.Places, seg.DestinationStation)
		segResponse := models.Segment{
			Origin:      models.Place{Name: origin.Name, Code: origin.Code},
			Destination: models.Place{Name: destination.Name, Code: destination.Code},
			Departure:   seg.DepartureDateTime,
			Arrival:     seg.ArrivalDateTime,
			Duration:    seg.Duration}
		segs = append(segs, segResponse)
	}

	return models.Leg{
		Departure:   leg.Departure,
		Arrival:     leg.Arrival,
		Duration:    leg.Duration,
		Stops:       stops,
		Origin:      models.Place{Name: origin.Name, Code: origin.Code},
		Destination: models.Place{Name: destination.Name, Code: destination.Code},
		Carriers:    carriers, Segments: segs}
}
