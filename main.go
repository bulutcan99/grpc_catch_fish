package main

import (
	"context"
	"github.com/bulutcan99/grpc_weather/internal/fetch"
	config_mongodb "github.com/bulutcan99/grpc_weather/pkg/config/mongodb"
	config_rabbitmq "github.com/bulutcan99/grpc_weather/pkg/config/rabbitmq"
	config_redis "github.com/bulutcan99/grpc_weather/pkg/config/redis"
	"github.com/bulutcan99/grpc_weather/pkg/env"
	"github.com/bulutcan99/grpc_weather/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	Mongo    *config_mongodb.Mongo
	Redis    *config_redis.Redis
	RabbitMQ *config_rabbitmq.RabbitMQ
	Logger   *zap.Logger
	Env      *env.ENV
)

func init() {
	Env = env.ParseEnv()
	Logger = logger.InitLogger(Env.LogLevel)
	Mongo = config_mongodb.NewConnetion()
	Redis = config_redis.NewRedisConnection()
	RabbitMQ = config_rabbitmq.NewRabbitMQConnection()
}

func main() {
	defer Logger.Sync()
	defer Mongo.Close()
	defer Redis.Close()
	defer RabbitMQ.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	data, err := fetch.FetchDefaultWeather()
	if err != nil {
		zap.S().Error("Error while fetching data: ", err)
	}
	zap.S().Info("Data: ", data)
	var wg sync.WaitGroup
	wg.Add(1)

	server := &http.Server{
		Addr: ":8081",
	}

	go func() {
		defer wg.Done()
		zap.S().Info("Server listening on :8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Fatal("HTTP server hatası: ", err)
		}
	}()

	<-ctx.Done()
	zap.S().Info("Shutting down server...")
	err = server.Shutdown(ctx)
	if err != nil {
		zap.S().Fatal("Sunucu kapatma hatası: ", err)
	}

	wg.Wait()
	zap.S().Info("Server gracefully stopped")
}
