package grpc_server

import (
	"context"
	"errors"
	"fmt"
	"github.com/bulutcan99/grpc_weather/model"
	"github.com/bulutcan99/grpc_weather/pkg/env"
	"github.com/bulutcan99/grpc_weather/proto/pb"
	"github.com/bulutcan99/grpc_weather/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

var (
	DEFAULT_CITY = &env.Env.DefaultWeatherCity
)

type WeatherServer struct {
	pb.UnimplementedUserServiceServer
	pb.UnimplementedWeatherServiceServer
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

		msg := fmt.Sprintf("Successfully logged in! User: %s", user)
		return &pb.ResponseLogin{
			Status:  msg,
			Success: true,
		}, nil
	}
}

func (s *WeatherServer) UpdatePassword(ctx context.Context, req *pb.RequestUpdate) (*pb.ResponseUpdate, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, errors.New("operation cancelled due to timeout")
	default:
		userId, err := primitive.ObjectIDFromHex(req.Id)
		if err != nil {
			return &pb.ResponseUpdate{
				Message: "User is not updated",
				Success: false,
			}, err
		}

		pass := strings.TrimSpace(req.Password)
		updatedUser, err := s.Services.UserService.UpdateUserPassword(userId, pass)
		if err != nil {
			return &pb.ResponseUpdate{
				Message: "User is not updated",
				Success: false,
			}, err
		}

		if updatedUser == nil {
			return &pb.ResponseUpdate{
				Message: "User is not updated",
				Success: false,
			}, errors.New("user is not updated")
		}

		msg := fmt.Sprintf("Successfully updated! User: %s", updatedUser.Username)
		return &pb.ResponseUpdate{
			Message: msg,
			Success: true,
		}, nil
	}
}

func (s *WeatherServer) GetUserCity(ctx context.Context, req *pb.RequstUserCity) (res *pb.ResponseUserCity, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, errors.New("operation cancelled due to timeout")
	default:
		city, err := s.Services.WeatherService.GetCityDataByUsername(req.Username)
		if err != nil {
			return &pb.ResponseUserCity{
				City:    "Error",
				Success: false,
			}, err
		}

		if city == "" {
			return &pb.ResponseUserCity{
				City:    "Error",
				Success: false,
			}, errors.New("city is empty")
		}

		return &pb.ResponseUserCity{
			City:    city,
			Success: true,
		}, nil
	}
}

func (s *WeatherServer) GetWeatherData(ctx context.Context, req *pb.RequestWeatherData) (*pb.ResponseWeatherData, error) {
	ticker := time.NewTicker(500 * time.Millisecond)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	weatherChan := make(chan *pb.ResponseWeatherData)
	errChan := make(chan error)
	go func() {
		city, err := s.Services.WeatherService.GetCityDataByUsername(req.Username)
		if city == "" || err != nil {
			city = *DEFAULT_CITY
		}

		city = strings.TrimSpace(city)
		weatherData, err := s.Services.WeatherService.FetchWeatherData(city)
		if err != nil {
			errChan <- errors.New("weather data is empty")
		}

		weather := &pb.Weather{
			City:        weatherData.City,
			Country:     weatherData.Country,
			Temperature: weatherData.TempC,
			CityTime:    weatherData.CityTime,
		}

		if weatherData.City == "" {
			errChan <- errors.New("weather data is empty")
			weatherChan <- &pb.ResponseWeatherData{
				Message: "Weather data is not fetched",
				Success: false,
			}
		}

		msg := fmt.Sprintf("Successfully fetched! Weather Temp: %v", weatherData.TempC)
		weatherChan <- &pb.ResponseWeatherData{
			Weather: weather,
			Message: msg,
			Success: true,
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("operation cancelled due to timeout")
		case <-ticker.C:
			zap.S().Info("Fetching data...")
		case err := <-errChan:
			return &pb.ResponseWeatherData{
				Message: "Error while fetching data",
				Success: false,
			}, err
		case weather := <-weatherChan:
			return weather, nil

		}
	}
}
