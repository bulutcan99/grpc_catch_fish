package model

type WeatherData struct {
	TempC    float64 `json:"temp_c" bson:"temp_c"`
	Country  string  `json:"country" bson:"country"`
	City     string  `json:"city" bson:"city"`
	CityTime string  `json:"city_time" bson:"city_time"`
}
