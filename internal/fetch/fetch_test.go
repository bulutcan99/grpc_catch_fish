package fetch

import (
	"errors"
	"fmt"
	"testing"
)

func TestFetch_Weather(t *testing.T) {
	client := NewFetchingDataClient("Istanbul")
	fmt.Println("Cl: ", client)
	data, err := client.FetchWeather(client.Url)
	if err != nil {
		errors.New("Error while getting data!")
	}
	fmt.Println("Data:", &data, data)
}
