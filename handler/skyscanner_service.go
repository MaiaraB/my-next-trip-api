package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/MaiaraB/travel-plan/ssmodels"
	"github.com/matryer/runner"
)

const (
	skyscannerPOSTFlightsSearchURL = "https://skyscanner-skyscanner-flight-search-v1.p.rapidapi.com/apiservices/pricing/v1.0"
	skyscannerGETFlightsSearchURL  = "https://skyscanner-skyscanner-flight-search-v1.p.rapidapi.com/apiservices/pricing/uk2/v1.0"
	skyscannerGETURL               = "https://skyscanner-skyscanner-flight-search-v1.p.rapidapi.com/apiservices/reference/v1.0"
	skyscannerSearchPlaceURL       = "https://skyscanner-skyscanner-flight-search-v1.p.rapidapi.com/apiservices/autosuggest/v1.0"
	skyscannerHost                 = "skyscanner-skyscanner-flight-search-v1.p.rapidapi.com"
	responseURL                    = "http://partners.api.skyscanner.net/apiservices/pricing/uk2/v1.0"
)

var skyscannerKey = os.Getenv("SKYSCANNER_API_KEY")

func getCountriesSkyscanner(locale string) []ssmodels.Country {
	client := &http.Client{Timeout: 10 * time.Second}

	request, _ := http.NewRequest("GET", skyscannerGETURL+"/countries/"+locale, nil)
	request.Header.Set("X-RapidAPI-Host", skyscannerHost)
	request.Header.Set("X-RapidAPI-Key", skyscannerKey)
	response, err := client.Do(request)

	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}

	responseData := ssmodels.CountryResponse{}

	if decodeErr := json.NewDecoder(response.Body).Decode(&responseData); decodeErr != nil {
		log.Printf("Decoding json failed with error %s\n", decodeErr)
	}

	return responseData.Countries
}

func getCurrenciesSkyscanner() []ssmodels.Currency {
	client := &http.Client{Timeout: 10 * time.Second}

	request, _ := http.NewRequest("GET", skyscannerGETURL+"/currencies", nil)
	request.Header.Set("X-RapidAPI-Host", skyscannerHost)
	request.Header.Set("X-RapidAPI-Key", skyscannerKey)
	response, err := client.Do(request)

	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}

	responseData := ssmodels.CurrencyResponse{}

	if decodeErr := json.NewDecoder(response.Body).Decode(&responseData); decodeErr != nil {
		log.Printf("Decoding json failed with error %s\n", decodeErr)
	}

	return responseData.Currencies
}

func getPlacesSkyscanner(country string, currency string, locale string, query string) []ssmodels.SearchPlace {
	client := &http.Client{Timeout: 10 * time.Second}
	url := skyscannerSearchPlaceURL + "/" + country + "/" + currency + "/" + locale + "/?query=" + query

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("X-RapidAPI-Host", skyscannerHost)
	request.Header.Set("X-RapidAPI-Key", skyscannerKey)
	response, err := client.Do(request)

	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}

	responseData := ssmodels.SearchPlaceResponse{}

	if decodeErr := json.NewDecoder(response.Body).Decode(&responseData); decodeErr != nil {
		log.Printf("Decoding json failed with error %s\n", decodeErr)
	}

	return responseData.Places
}

func getFlightsSkyscanner(shouldStop runner.S, respond chan<- ssmodels.FlightsResponse, data url.Values, index int) {
	client := &http.Client{Timeout: 10 * time.Second}

	response := postFlightSearch(data, client, index)

	if shouldStop() {
		return
	}

	for response.StatusCode == 429 {
		log.Printf("Thread #%d: TOO MANY POST REQUESTS", index)
		time.Sleep(time.Second)

		response = postFlightSearch(data, client, index)
		if shouldStop() {
			return
		}
	}

	locationKey := response.Header.Get("Location")[len(responseURL):]

	flightResponse := getFlightSearchResults(locationKey, client, index)

	if shouldStop() {
		return
	}

	log.Printf("FLIGHT RESPONSE #%d: %+v. Size: %d", index, flightResponse.Query, len(flightResponse.Itineraries))

	respond <- flightResponse
}

func postFlightSearch(data url.Values, client *http.Client, index int) *http.Response {
	request, _ := http.NewRequest("POST", skyscannerPOSTFlightsSearchURL, strings.NewReader(data.Encode()))
	request.Header.Set("X-RapidAPI-Host", skyscannerHost)
	request.Header.Set("X-RapidAPI-Key", skyscannerKey)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	printRequestBody(request, index)

	response, err := client.Do(request)
	// log.Println("POST RESPONSE: ", response)

	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		panic(err)
	}

	log.Printf("POST #%d response header: %s\nPOST response status: %s", index, response.Header, response.Status)

	return response
}

func getFlightSearchResults(locationKey string, client *http.Client, index int) ssmodels.FlightsResponse {
	request, _ := http.NewRequest("GET", skyscannerGETFlightsSearchURL+locationKey+"?sortType=price&sortOrder=asc&pageIndex=0&pageSize=20", nil)
	request.Header.Set("X-RapidAPI-Host", skyscannerHost)
	request.Header.Set("X-RapidAPI-Key", skyscannerKey)
	response, err := client.Do(request)

	responseData := ssmodels.FlightsResponse{}

	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	} else {
		if decodeErr := json.NewDecoder(response.Body).Decode(&responseData); decodeErr != nil {
			log.Printf("Decoding json failed with error %s\n", decodeErr)
		}
		// log.Printf("FIRST GET RESPONSE #%d: %v", index, responseData)
		for response.StatusCode == 429 || responseData.Status == "UpdatesPending" {
			time.Sleep(time.Second)
			request, _ = http.NewRequest("GET", skyscannerGETFlightsSearchURL+locationKey+"?sortType=price&sortOrder=asc&pageIndex=0&pageSize=20", nil)
			request.Header.Set("X-RapidAPI-Host", skyscannerHost)
			request.Header.Set("X-RapidAPI-Key", skyscannerKey)
			response, err = client.Do(request)
			// log.Printf("PENDING GET RESPONSE #%d: %v", index, response.StatusCode)
			if err != nil {
				log.Printf("The HTTP request failed with error %s\n", err)
			} else {
				responseData = ssmodels.FlightsResponse{}
				if decodeErr := json.NewDecoder(response.Body).Decode(&responseData); decodeErr != nil {
					log.Printf("Decoding json failed with error %s\n", decodeErr)
				}
				// log.Printf("PENDING GET STATUS #%d: %v", index, responseData.Status)
			}
		}
	}

	return responseData
}

func printRequestBody(request *http.Request, index int) {
	buf, bodyErr := ioutil.ReadAll(request.Body)
	if bodyErr != nil {
		log.Print("bodyErr ", bodyErr.Error())
		// http.Error(w, bodyErr.Error(), http.StatusInternalServerError)
		// return
	}

	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
	log.Printf("BODY #%d: %q", index, rdr1)
	request.Body = rdr2
}
