package dto

type CurrentWeather struct {
	Cloud       int     `json:"cloud"`
	IsDay       int     `json:"is_day"`
	LastUpdated string  `json:"last_updated"`
	TempC       float64 `json:"temp_c"`
	UV          float64 `json:"uv"`
	WindDegree  int     `json:"wind_degree"`
	WindDir     string  `json:"wind_dir"`
	WindKph     float64 `json:"wind_kph"`
}

type LocationData struct {
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Localtime      string  `json:"localtime"`
	LocaltimeEpoch float64 `json:"localtime_epoch"`
	Lon            float64 `json:"lon"`
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	TzID           string  `json:"tz_id"`
}

type WeatherData struct {
	Current  CurrentWeather `json:"current"`
	Location LocationData   `json:"location"`
}
