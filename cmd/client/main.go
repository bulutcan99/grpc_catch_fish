package main

import (
	"github.com/bulutcan99/grpc_weather/api/grpc_client"
	"go.uber.org/zap"
)

func main() {
	weatherClient := grpc_client.NewWeatherClient()
	defer weatherClient.Close()
	if err := weatherClient.GetWeatherDataByLatLong(); err != nil {
		zap.S().Fatal(err)
	}
}
