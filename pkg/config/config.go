package config_builder

import (
	"github.com/bulutcan99/grpc_weather/pkg/env"
	"github.com/gofiber/fiber/v2"
	"time"
)

var READ_TIMEOUT_SECONDS_COUNT = &env.Env.ServerReadTimeout

func ConfigFiber() fiber.Config {
	return fiber.Config{
		ReadTimeout:  time.Second * time.Duration(*READ_TIMEOUT_SECONDS_COUNT),
		Prefork:      false,
		ServerHeader: "iMon",
		AppName:      "Go-Chat",
		Immutable:    true,
	}
}
