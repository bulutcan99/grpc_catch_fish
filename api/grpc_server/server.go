package grpc_server

import (
	"context"
	"errors"
	"fmt"
	"github.com/bulutcan99/grpc_weather/model"
	"github.com/bulutcan99/grpc_weather/service"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	defaultCity = "Istanbul"
)

type WeatherServer struct {
	pb.UnimplementedUserServiceServer
	mutex    *sync.Mutex
	Services *service.Services
	cityConn chan struct{}
}

func NewWeatherServer(Services *service.Services) *WeatherServer {
	return &WeatherServer{
		Services: Services,
		mutex:    new(sync.Mutex),
		cityConn: make(chan struct{}),
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

		userId, err := s.Services.UserService.RegisterUser(user)
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

func (s *WeatherServer) Login(ctx context.Context, req *pb.RequestLogin) (*pb.ResponseLogin, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, errors.New("operation cancelled due to timeout")
	default:
		user, err := s.Services.UserService.LoginUser(req.Username, req.Password)
		if err != nil {
			return &pb.ResponseLogin{
				Status:  "User is not found",
				Success: false,
			}, err
		}

		if user == nil {
			return &pb.ResponseLogin{
				Status:  "User is not found",
				Success: false,
			}, errors.New("user is not found")
		}

		msg := fmt.Sprintf("Successfully logged in! User: %s", user.Username)
		return &pb.ResponseLogin{
			Status:  msg,
			Success: true,
		}, nil
	}
}
