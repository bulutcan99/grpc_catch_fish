package main

import (
	"context"
	"encoding/json"
	config_mongodb "github.com/bulutcan99/grpc_weather/pkg/config/mongodb"
	config_rabbitmq "github.com/bulutcan99/grpc_weather/pkg/config/rabbitmq"
	config_redis "github.com/bulutcan99/grpc_weather/pkg/config/redis"
	"github.com/bulutcan99/grpc_weather/pkg/env"
	"github.com/bulutcan99/grpc_weather/pkg/logger"
	"go.uber.org/zap"
	"io"
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

	apiURL := "http://api.weatherapi.com/v1/current.json?key=5e991bddf944431e858131733232111&q=London&aqi=no"

	response, err := http.Get(apiURL)
	if err != nil {
		zap.S().Error("API Request failed: ", err)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		zap.S().Error("API response failed: ", err)
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		zap.S().Error("There is an error while parsing data: ", err)
		return
	}

	zap.S().Info("Data: ", data)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	server := &http.Server{
		Addr: ":8080",
	}

	http.HandleFunc("/api/example/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	})

	go func() {
		defer wg.Done()
		zap.S().Info("Server listening on :8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Fatal("HTTP server hatası: ", err)
		}
	}()

	<-ctx.Done()

	zap.S().Info("Shutting down server...")

	// Sunucuyu kapat
	err = server.Shutdown(context.Background())
	if err != nil {
		zap.S().Fatal("Sunucu kapatma hatası: ", err)
	}

	// WaitGroup bekleme sürecini tamamla
	wg.Wait()

	zap.S().Info("Server gracefully stopped")
}
