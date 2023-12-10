package service

import (
	"errors"
	"github.com/bulutcan99/grpc_weather/internal/fetch"
	"github.com/bulutcan99/grpc_weather/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IWeatherRepo interface {
	InsertOrUpdate(filter any, update any) *mongo.SingleResult
	FindOne(filter any) *mongo.SingleResult
}

type WeatherService struct {
	WeatherRepo IWeatherRepo
	UserRepo    IUserRepo
	Fetcher     *fetch.FetchingDataClient
}

func NewWeatherService(weatherRepo IWeatherRepo, userRepo IUserRepo) *WeatherService {
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

	return weatherData, nil
}

func (w *WeatherService) UpdateWeatherData(weatherData *model.WeatherData) error {
	filter := bson.D{
		{"city", weatherData.City},
	}

	update := bson.D{
		{"temp_c", weatherData.TempC},
		{"country", weatherData.Country},
		{"city_time", weatherData.CityTime},
	}

	response := w.WeatherRepo.InsertOrUpdate(filter, update)
	if err := response.Err(); err != nil {
		return errors.New("error while updating data")
	}

	return nil
}

func (w *WeatherService) GetWeatherData(city string) (*model.WeatherData, error) {
	var data *model.WeatherData
	filter := bson.D{
		{"city", city},
	}

	response := w.WeatherRepo.FindOne(filter)
	if err := response.Decode(&data); err != nil {
		return nil, errors.New("error while fetching data")
	}

	return data, nil
}

func (w *WeatherService) GetCityDataByUsername(username string) (string, error) {
	filter := bson.D{
		{"username", username},
	}

	opt := options.FindOne().SetProjection(bson.D{{"city", 1}})
	response, err := w.UserRepo.FindOne(filter, opt)
	if err != nil {
		return "Error", errors.New("error while finding user")
	}

	city := response.City
	return city, nil
}

func (w *WeatherService) GetUserWeatherData(username string) (*model.WeatherData, error) {
	city, err := w.GetCityDataByUsername(username)
	if err != nil {
		return nil, err
	}

	weatherData, err := w.GetWeatherData(city)
	if err != nil {
		return nil, err
	}

	return weatherData, nil
}
