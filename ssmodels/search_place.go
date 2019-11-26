package ssmodels

type SearchPlace struct {
	PlaceID     string `json:"PlaceId"`
	PlaceName   string
	CountryID   string `json:"CountryId"`
	RegionID    string `json:"RegionId"`
	CityID      string `json:"CityId"`
	CountryName string
}
