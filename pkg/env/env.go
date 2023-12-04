package env

import (
	"fmt"
	custom_error "github.com/bulutcan99/grpc_weather/pkg/error"
	"os"
	"sync"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type ENV struct {
	ServerHost         string `env:"SERVER_HOST,required"`
	ServerPort         int    `env:"SERVER_PORT,required"`
	DbPort             int    `env:"DB_PORT,required"`
	DbName             string `env:"DB_NAME,required"`
	UserCollection     string `env:"USER_COLLECTION,required"`
	WeatherUrl         string `env:"WEATHER_URL,required"`
	WeatherApiKey      string `env:"WEATHER_API_KEY,required"`
	DefaultWeatherCity string `env:"DEFAULT_WEATHER_CITY,required"`
	RabbitMQPort       int    `env:"RABBITMQ_PORT,required"`
	RabbitMQUser       string `env:"RABBITMQ_USER,required"`
	RabbitMQPassword   string `env:"RABBITMQ_PASSWORD,required"`
	RedisPort          int    `env:"REDIS_PORT,required"`
	RedisPassword      string `env:"REDIS_PASSWORD,required"`
	RedisDBNumber      int    `env:"REDIS_DB_NUMBER,required"`
	LogLevel           string `env:"LOG_LEVEL,required"`
}

var doOnce sync.Once
var Env ENV

func ParseEnv() *ENV {
	doOnce.Do(func() {
		e := godotenv.Load()
		if e != nil {
			custom_error.ParseError()
			os.Exit(1)
		}
		if err := env.Parse(&Env); err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(0)
		}
	})
	return &Env
}
