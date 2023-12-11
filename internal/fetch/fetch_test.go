package fetch

import (
	"errors"
	"fmt"
	"testing"
)

func TestFetch_Weather(t *testing.T) {
	client := NewFetchingDataClient()
	data, err := client.FetchWeather("http://api.weatherapi.com/v1/current.json?key=&q=London&aqi=no")
	if err != nil {
		errors.New("Error while getting data!")
	}
	fmt.Println("Data:", data)
}
