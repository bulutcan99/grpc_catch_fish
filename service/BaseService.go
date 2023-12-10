package service

import (
	"github.com/bulutcan99/grpc_weather/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"

	"go.uber.org/zap"
)

type IUserService interface {
	RegisterUser(newUser model.User) (primitive.ObjectID, error)
	LoginUser(username string, password string) (*model.User, error)
	UpdateUserPassword(id primitive.ObjectID, pass string) (*model.User, error)
}

type IWeatherService interface {
	FetchWeatherData(city string) (*model.WeatherData, error)
	UpdateWeatherData(weatherData *model.WeatherData) error
	GetWeatherData(city string) (*model.WeatherData, error)
	GetCityDataByUsername(username string) (string, error)
}

type Services struct {
	UserService    IUserService
	WeatherService IWeatherService
}

func RegisterServices(services ...any) *Services {

	newServices := &Services{}
	for _, service := range services {
		switch getServiceType(service) {
		case "UserService":
			newServices.UserService = service.(IUserService)
			zap.S().Info("UserService registered")
		case "WeatherService":
			newServices.WeatherService = service.(IWeatherService)
			zap.S().Info("WeatherService registered")
		}
	}
	return newServices
}

func getServiceType(service any) string {
	valueOf := reflect.ValueOf(service)
	var name string
	if valueOf.Type().Kind() == reflect.Ptr {
		name = reflect.Indirect(valueOf).Type().Name()
	} else {
		name = valueOf.Type().Name()
	}
	return name
}
