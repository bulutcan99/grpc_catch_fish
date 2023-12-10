package fetch

import (
	"errors"
	"fmt"
	"testing"
)

func TestFetch_Weather(t *testing.T) {
	client := NewFetchingDataClient()
	data, err := client.FetchWeather("https://api.weatherapi.com/v1/current.json?q=Istanbul&lang=en&key=290ea74ca97c4725824222051230912")
	if err != nil {
		errors.New("Error while getting data!")
	}
	fmt.Println("Data:", data)
}
