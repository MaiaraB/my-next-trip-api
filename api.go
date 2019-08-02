/*
 * Flight Planner
 *
 * This is a Flight Planner server.
 *
 * API version: 1.0.0
 * Contact: maiarabarroso84@gmail.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	layoutISO = "2006-01-02"
)

func getFlights(w http.ResponseWriter, r *http.Request) {
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		log.Println("The client closed the connection prematurely. Cleaning up.")
		// panic(http.ErrAbortHandler)
	}()

	queryValues := r.URL.Query()
	// roundTrip, _ := strconv.ParseBool(queryValues.Get("roundTrip"))
	origin := queryValues.Get("origin")
	destination := queryValues.Get("destination")
	outboundWeekDay, _ := strconv.Atoi(queryValues.Get("outboundWeekDay"))
	duration, _ := strconv.Atoi(queryValues.Get("duration"))
	country := queryValues.Get("country")
	currency := queryValues.Get("currency")
	locale := queryValues.Get("locale")

	intervals := getDateIntervals(time.Weekday(outboundWeekDay), duration, time.Now(), time.Now().AddDate(0, 3, 0))
	data := url.Values{}
	data.Set("cabinClass", "economy")
	data.Set("adults", "1")
	data.Set("locale", locale)
	data.Set("currency", currency)
	data.Set("country", country)
	data.Set("originPlace", origin)
	data.Set("destinationPlace", destination)

	// var counter = 0

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, HEAD, GET")
	// w.Header().Set("Connection", "Keep-Alive")
	// w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Printf("Expected http.ResponseWriter to be an http.Flusher")

	}
	w.Write([]byte(fmt.Sprintf("%d<", len(intervals))))
	flusher.Flush()

	results := make(chan []FlightsResult, len(intervals))

	for i := 0; i < len(intervals); i++ {
		dataCopy := copyValues(data)
		dataCopy.Set("outboundDate", intervals[i].Outbound.Format(layoutISO))
		dataCopy.Set("inboundDate", intervals[i].Inbound.Format(layoutISO))

		go getPartialResults(results, dataCopy)
	}

	for i := 0; i < len(intervals); i++ {
		// json.NewEncoder(w).Encode(flightResults)
		flightResultJSON, err := json.Marshal(<-results)
		if err != nil {
			fmt.Printf("The response marshalling failed with error %s\n", err)
		}

		w.Write(flightResultJSON)
		w.Write([]byte("<"))
		flusher.Flush()
	}
	//log.Println("#results: ", counter)

}

func copyValues(data url.Values) url.Values {
	newData := url.Values{}
	for k, v := range data {
		newData[k] = v
	}
	return newData
}

func getPartialResults(respond chan<- []FlightsResult, data url.Values) {
	flightResponse := getItineraries(data)
	var flightResults []FlightsResult

	for _, it := range flightResponse.Itineraries {
		result := FlightsResult{}

		// Setting result Currency
		result.Currency = searchCurrencyByCode(flightResponse.Currencies, flightResponse.Query.Currency)

		// Setting result AgentInfo
		var agentsInfo []AgentInfo
		for _, po := range it.PricingOptions {
			agent := searchAgentByID(flightResponse.Agents, po.AgentIds[0])
			agentsInfo = append(agentsInfo, AgentInfo{agent.Name, agent.ImageURL, po.Price, po.DeeplinkURL})
		}
		result.AgentsInfo = agentsInfo

		// Setting result InboundLeg
		result.InboundLeg = configResponseLeg(flightResponse, it.InboundLegID)

		// Setting result OutboundLeg
		result.OutboundLeg = configResponseLeg(flightResponse, it.OutboundLegID)

		flightResults = append(flightResults, result)

		// log.Println("RESULT: ", result)
		// counter++
		// log.Println("Itinerary #", counter)
		log.Println("Itinerary price: ", result.AgentsInfo[0].Price)
	}

	respond <- flightResults
}

func configResponseLeg(flightResponse FlightsResponse, legID string) LegResponse {

	leg := searchLegByID(flightResponse.Legs, legID)

	stops := []PlaceResponse{}
	for _, stopID := range leg.StopsIds {
		stop := searchPlaceByID(flightResponse.Places, stopID)
		stops = append(stops, PlaceResponse{stop.Name, stop.Code})
	}

	origin := searchPlaceByID(flightResponse.Places, leg.OriginStation)
	destination := searchPlaceByID(flightResponse.Places, leg.DestinationStation)

	carriers := []CarrierResponse{}
	for _, carrierID := range leg.CarrierIds {
		carrier := searchCarrierByID(flightResponse.Carriers, carrierID)
		carriers = append(carriers, CarrierResponse{carrier.Name, carrier.ImageURL})
	}

	segs := []SegmentResponse{}
	for _, segID := range leg.SegmentIds {
		seg := searchSegmentByID(flightResponse.Segments, segID)
		origin := searchPlaceByID(flightResponse.Places, seg.OriginStation)
		destination := searchPlaceByID(flightResponse.Places, seg.DestinationStation)
		segResponse := SegmentResponse{
			PlaceResponse{origin.Name, origin.Code},
			PlaceResponse{destination.Name, destination.Code},
			seg.DepartureDateTime,
			seg.ArrivalDateTime,
			seg.Duration}
		segs = append(segs, segResponse)
	}

	return LegResponse{leg.Departure, leg.Arrival, leg.Duration, stops,
		PlaceResponse{origin.Name, origin.Code},
		PlaceResponse{destination.Name, destination.Code},
		carriers, segs}
}

type FlightsResult struct {
	Currency    Currency
	AgentsInfo  []AgentInfo
	InboundLeg  LegResponse
	OutboundLeg LegResponse
}

type AgentInfo struct {
	Name        string
	ImageURL    string
	Price       float64
	DeepLinkURL string
}

type LegResponse struct {
	Departure   MyTime
	Arrival     MyTime
	Duration    int
	Stops       []PlaceResponse // place codes
	Origin      PlaceResponse   // place code
	Destination PlaceResponse   //place code
	Carriers    []CarrierResponse
	Segments    []SegmentResponse
}

type CarrierResponse struct {
	Name     string
	ImageURL string
}

type SegmentResponse struct {
	Origin      PlaceResponse
	Destination PlaceResponse
	Departure   MyTime
	Arrival     MyTime
	Duration    int
}

type PlaceResponse struct {
	Name string
	Code string
}
