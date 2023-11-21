package config_builder

import (
	"fmt"
	"github.com/bulutcan99/grpc_weather/pkg/env"
)

var (
	SEVER_HOST        = &env.Env.ServerHost
	DB_PORT           = &env.Env.DbPort
	REDIS_PORT        = &env.Env.RedisPort
	GRPC_PORT         = &env.Env.GrpcPort
	RABBITMQ_USER     = &env.Env.RabbitMQUser
	RABBITMQ_PASSWORD = &env.Env.RabbitMQPassword
	RABBITMQ_PORT     = &env.Env.RabbitMQPort
)

func ConnectionURLBuilder(n string) (string, error) {
	var url string
	switch n {
	case "grpc":
		url = fmt.Sprintf(
			"%s:%d",
			*SEVER_HOST,
			*GRPC_PORT,
		)
	case "mongo":
		url = fmt.Sprintf(
			"mongodb://%s:%d",
			*SEVER_HOST,
			*DB_PORT,
		)
	case "rabbitmq":
		url = fmt.Sprintf("amqp://%s:%s@%s:%d/",
			*RABBITMQ_USER,
			*RABBITMQ_PASSWORD,
			*SEVER_HOST,
			*RABBITMQ_PORT,
		)
	case "redis":
		url = fmt.Sprintf(
			"%s:%d",
			*SEVER_HOST,
			*REDIS_PORT,
		)
	default:
		return "", fmt.Errorf("connection name '%v' is not supported", n)
	}

	return url, nil
}
