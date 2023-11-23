package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bulutcan99/grpc_weather/dto"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const (
	apiURL      = "http://api.weatherapi.com/v1/current.json"
	apiKey      = "5e991bddf944431e858131733232111"
	defaultCity = "London"
)

func buildURL(city string) string {
	return fmt.Sprintf("%s?key=%s&q=%s&aqi=no", apiURL, apiKey, city)
}

func fetchData(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func parseData(data []byte) (dto.WeatherData, error) {
	var result dto.WeatherData
	err := json.Unmarshal(data, &result)
	if err != nil {
		return dto.WeatherData{}, err
	}

	return result, nil
}

func FetchWeather(city string) (dto.WeatherData, error) {
	url := buildURL(city)

	body, err := fetchData(url)
	if err != nil {
		zap.S().Error("Failed to fetch data: ", err)
		return dto.WeatherData{}, errors.New("failed to fetch data")
	}

	data, err := parseData(body)
	if err != nil {
		zap.S().Error("Failed to parse data: ", err)
		return dto.WeatherData{}, errors.New("failed to parse data")
	}

	return data, nil
}

func FetchDefaultWeather() (dto.WeatherData, error) {
	return FetchWeather(defaultCity)
}
