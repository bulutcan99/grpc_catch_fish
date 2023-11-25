package main

import (
	"context"
	"fmt"
	"github.com/bulutcan99/grpc_weather/api/grpc_server"
	"github.com/bulutcan99/grpc_weather/internal/query"
	config_mongodb "github.com/bulutcan99/grpc_weather/pkg/config/mongodb"
	config_rabbitmq "github.com/bulutcan99/grpc_weather/pkg/config/rabbitmq"
	config_redis "github.com/bulutcan99/grpc_weather/pkg/config/redis"
	"github.com/bulutcan99/grpc_weather/pkg/env"
	"github.com/bulutcan99/grpc_weather/pkg/logger"
	"github.com/bulutcan99/grpc_weather/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
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

func Init() {
	Env = env.ParseEnv()
	Logger = logger.InitLogger(Env.LogLevel)
	Mongo = config_mongodb.NewConnetion()
	Redis = config_redis.NewRedisConnection()
	RabbitMQ = config_rabbitmq.NewRabbitMQConnection()
}

func main() {
	Init()
	defer Logger.Sync()
	defer Mongo.Close()
	defer Redis.Close()
	defer RabbitMQ.Close()
	fmt.Println("----------------------")
	fmt.Println("3131")
	userRepo := query.NewUserRepositry(Mongo, Env.UserCollection)
	userService := service.NewUserService(userRepo)
	grpcServer := grpc.NewServer()
	grpc_server.NewWeatherServer(userService, grpcServer)
	fmt.Println("gRPC Weather Server starting on: ", Env.GrpcPort)
	lis, err := net.Listen("tcp", Env.GrpcPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	<-ctx.Done()

	zap.S().Info("Shutting down server...")
	grpcServer.GracefulStop()

	wg.Wait()
	zap.S().Info("Server gracefully stopped")
}
