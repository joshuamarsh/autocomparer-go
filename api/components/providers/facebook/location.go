package facebook

// LocationResponse response parser for mapquest
type LocationResponse struct {
	Status int `json:"status"`
	Result struct {
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
	} `json:"result"`
}

// LocationCityResponse response parser for mapquest
type LocationCityResponse struct {
	Results []struct {
		Locations []struct {
			LatLng struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"latLng"`
		} `json:"locations"`
	} `json:"results"`
}
