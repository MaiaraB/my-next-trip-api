# MyNextTripAPI

An API built in Go and used by [this application](https://github.com/MaiaraB/my-next-trip), mainly for getting the results for flights searches.

The available API calls are:

## getFlights

Endpoint that returns all the flights from `origin` to `destionation` departuring on the giving `outboundWeekDay` and with the giving `duration` in the date range between `fromDate` and `toDate`.

### URL

http://localhost:<env PORT>/api/flights
  
### Method

`GET`

### URL Params

Required:
`cabinClass=[economy | premiumeconomy | business | first]`
`origin=[placeID]`
`destination=[placeID]`
`outboundWeekDay=[0-6]`
`adults=[1-8]`
`country=[ex: US]`
`locale=[ISO locale]`
`currency=[3-letter currency code]`
`fromDate=[YYYY-MM-DD]`
`toDate=[YYYY-MM-DD]`

Optional:
`duration=[1-30]` (not necessary for one way trips)
`children=[0-8]`
`infants=[0-8]`

### Success response

Code: 200
Content: `[# of chunks]<[Chunk1]<...<[ChunkN]`
Chunk content: 
  `[
    {Currency: {Code:, Symbol:, ThousandsSeparator:, DecimalSeparator:, SymbolOnLeft:, SpaceBetweenAmountAndSymbol:, RoundingCoefficient:, DecimalDigits:}, 
     AgentsInfo: [{Name:, ImageURL:, Price:, DeepLinkURL: }], 
     InboundLeg: {Departure:, Arrival:, Duration:, Stops:[], Origin: {Name:, Code:}, Destination: {Name:, Code:}, Carriers: [{Name:, ImageURL:}], Segments: [{Origin:, Destination:, Departure:, Arrival:}]}, 
     OutboundLeg: {similar to InboundLeg}
     }, ...
   ]`
   
## getCountries

### URL

http://localhost:<env PORT>/api/countries
  
### Method

`GET`

### URL Params

Required:
`locale=[ISO locale]`

### Success response

Code: 200
Content: `[{Code: , Name: },...]`

## getCurrencies

### URL

http://localhost:<env PORT>/api/currencies
  
### Method

`GET`

### Success response

Code: 200
Content: `[{Currency: {Code:, Symbol:, ThousandsSeparator:, DecimalSeparator:, SymbolOnLeft:, SpaceBetweenAmountAndSymbol:, RoundingCoefficient:, DecimalDigits:}, ...]`

## getPlaces

### URL

http://localhost:<env PORT>/api/queryPlace
  
### Method

`GET`

### URL Params

Required:
`country=[ex: US]`
`locale=[ISO locale]`
`currency=[3-letter currency code]`
`query=[string]`

### Success response

Code: 200
Content: `[{PlaceId: ,PlaceName: ,CountryId: ,RegionId: ,CityId: ,CountryName: }, ...]
