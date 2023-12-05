package service

import (
	"errors"
	"github.com/bulutcan99/grpc_weather/internal/fetch"
	"github.com/bulutcan99/grpc_weather/internal/query"
	"github.com/bulutcan99/grpc_weather/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IWeatherService interface {
	FetchWeatherData(city string) (*model.WeatherData, error)
	UpdateWeatherData(weatherData *model.WeatherData) error
	GetWeatherData(city string) (*model.WeatherData, error)
}

type WeatherService struct {
	WeatherRepo query.IWeatherRepo
	UserRepo    query.IUserRepo
	Fetcher     *fetch.FetchingDataClient
}

func NewWeatherService(weatherRepo query.IWeatherRepo, userRepo query.IUserRepo) *WeatherService {
	return &WeatherService{
		WeatherRepo: weatherRepo,
		UserRepo:    userRepo,
		Fetcher:     fetch.NewFetchingDataClient(),
	}
}

func (w *WeatherService) FetchWeatherData(city string) (*model.WeatherData, error) {
	url := w.Fetcher.GetURL(city)
	weatherData, err := w.Fetcher.FetchWeather(url)
	if err != nil {
		return nil, err
	}

	err = w.UpdateWeatherData(weatherData)
	if err != nil {
		return nil, err
	}

	return w.GetWeatherData(weatherData.City)
}

func (w *WeatherService) UpdateWeatherData(weatherData *model.WeatherData) error {
	filter := bson.D{
		{"city", weatherData.City},
	}

	update := bson.D{
		{"$set", bson.D{
			{"temp_c", weatherData.TempC},
			{"country", weatherData.Country},
			{"city_time", weatherData.CityTime},
		}},
	}

	response := w.WeatherRepo.InsertOrUpdate(filter, update)
	if err := response.Err(); err != nil {
		return errors.New("error while updating data")
	}

	return nil
}

func (w *WeatherService) GetWeatherData(city string) (*model.WeatherData, error) {
	var data model.WeatherData
	filter := bson.D{
		{"city", city},
	}

	response := w.WeatherRepo.FindOne(filter)
	if err := response.Decode(&data); err != nil {
		return nil, errors.New("error while fetching data")
	}

	return &data, nil
}

func (w *WeatherService) GetWeatherDataByUser(id primitive.ObjectID) (string, error) {
	filter := bson.D{
		{"_id", id},
	}

	response, err := w.UserRepo.FindOne(filter)
	if err != nil {
		return "Error", errors.New("error while finding user")
	}

	city := response.City
	return city, nil
}
