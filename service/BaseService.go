package service

import (
	"reflect"

	"go.uber.org/zap"
)

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
