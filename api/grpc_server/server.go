package grpc_server

import (
	"context"
	"errors"
	"github.com/bulutcan99/grpc_weather/model"
	pb "github.com/bulutcan99/grpc_weather/proto"
	"github.com/bulutcan99/grpc_weather/service"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	defaultCity = "Istanbul"
)

type WeatherServer struct {
	pb.UnimplementedWeatherServiceServer
	mutex       *sync.Mutex
	userService *service.UserService
	cityConn    chan struct{}
}

func NewWeatherServer(userService *service.UserService) *WeatherServer {
	return &WeatherServer{
		userService: userService,
		mutex:       new(sync.Mutex),
		cityConn:    make(chan struct{}),
	}
}

func (s *WeatherServer) Register(ctx context.Context, req *pb.RequestRegister) (*pb.ResponseRegister, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, errors.New("operation cancelled due to timeout")
	default:
		user := model.User{
			Username: req.Username,
			Password: req.Password,
			Name:     req.Name,
			Email:    req.Email,
			City:     req.City,
		}

		userId, err := s.userService.RegisterUser(user)
		if err != nil {
			return &pb.ResponseRegister{
				Message: "User is not registered",
				Success: false,
			}, err
		}

		zap.S().Info("User is registered with id: ", userId)
		return &pb.ResponseRegister{
			Message: "User is registered",
			Success: true,
		}, nil
	}
}
