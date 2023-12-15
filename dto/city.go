package dto

type Property struct {
	Name        string `json:"name"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	Postcode    string `json:"postcode"`
}

type Feature struct {
	Type       string   `json:"type"`
	Properties Property `json:"properties"`
}

type FetchedCityData struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}
