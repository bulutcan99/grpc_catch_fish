package fetch

import (
	"errors"
	"fmt"
	"github.com/bulutcan99/grpc_weather/dto"
	"github.com/bulutcan99/grpc_weather/model"
	config_http "github.com/bulutcan99/grpc_weather/pkg/config/http"
	"github.com/bulutcan99/grpc_weather/pkg/env"
	decoder "github.com/goccy/go-json"
)

var (
	API_URL = &env.Env.WeatherUrl
	API_KEY = &env.Env.WeatherApiKey
)

type FetchingDataClient struct {
	client *config_http.HttpClient
}

func NewFetchingDataClient() *FetchingDataClient {
	return &FetchingDataClient{
		client: config_http.NewHttpClient(),
	}
}

func (f *FetchingDataClient) GetURL(city string) string {
	url := fmt.Sprintf("%s?q=%s&lang=en&key=%s", *API_URL, city, *API_KEY)
	return url
}

func (f *FetchingDataClient) parseData(data []byte) (*model.WeatherData, error) {
	var result dto.WeatherData
	if err := decoder.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	weatherData := &model.WeatherData{
		TempC:    result.Current.TempC,
		Country:  result.Location.Country,
		City:     result.Location.Name,
		CityTime: result.Location.Localtime,
	}
	return weatherData, nil
}

func (f *FetchingDataClient) FetchWeather(url string) (*model.WeatherData, error) {
	if f.client == nil {
		return nil, errors.New("http client is not initialized")
	}

	body, err := f.client.Get(url)
	if err != nil {
		return nil, err
	}

	data, err := f.parseData(body)
	if err != nil {
		return nil, err
	}

	return data, nil

}

// ticker := time.NewTicker(5 * time.Second)
// go func() {
// 	defer wg.Done()
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-ticker.C:
// 			zap.S().Info("Fetching data...")
// 			data, err := fetch.FetchWeather("izmir")
// 			if err != nil {
// 				zap.S().Error("Error while fetching data: ", err)
// 			}
// 			zap.S().Info("Data: ", data)
// 		}
// 	}
// }()
