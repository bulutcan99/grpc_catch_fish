package dto

type CurrentWeather struct {
	Cloud int     `json:"cloud"`
	TempC float64 `json:"temp_c"`
}

type LocationData struct {
	Country   string `json:"country"`
	Localtime string `json:"localtime"`
	Name      string `json:"name"`
}

type WeatherData struct {
	Current  CurrentWeather `json:"current"`
	Location LocationData   `json:"location"`
}
