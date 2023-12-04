package service

import (
	"reflect"

	"go.uber.org/zap"
)

type Services struct {
	UserService *UserService
	// WeatherService *WeatherService
}

func RegisterServices(services ...any) *Services {

	newServices := &Services{}

	for _, service := range services {
		switch getServiceType(service) {
		case "UserService":
			newServices.UserService = service.(*UserService)
			zap.S().Info("UserService registered")
			// case "WeatherService":
			// 	newServices.WeatherService = service.(*WeatherService)
			// 	zap.S().Info("WeatherService registered")
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
