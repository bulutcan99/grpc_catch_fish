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
	var wg sync.WaitGroup
	startTime := time.Now()
	ticker := time.NewTicker(250 * time.Millisecond)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	weatherChan := make(chan *pb.ResponseWeatherData)
	errChan := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()

		city, err := s.Services.WeatherService.GetCityDataByUsername(req.Username)
		if city == "" || err != nil {
			city = *DEFAULT_CITY
		}

		city = strings.TrimSpace(city)
		time.Sleep(2 * time.Second)
		weatherData, err := s.Services.WeatherService.FetchWeatherData(city)
		if err != nil {
			errChan <- errors.New("weather data is empty")
			return
		}

		weather := &pb.Weather{
			City:        weatherData.City,
			Country:     weatherData.Country,
			Temperature: weatherData.TempC,
			CityTime:    weatherData.CityTime,
		}

		if weatherData.City == "" {
			errChan <- errors.New("weather data is empty")
			return
		}

		msg := fmt.Sprintf("Successfully fetched! Weather Temp: %v", weatherData.TempC)
		weatherChan <- &pb.ResponseWeatherData{
			Weather: weather,
			Message: msg,
			Success: true,
		}
	}()

	go func() {
		wg.Wait()
		close(weatherChan)
		close(errChan)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("operation cancelled due to timeout")
		case <-ticker.C:
			elapsedTime := time.Since(startTime).Seconds()
			elapsedTimeFormatted := fmt.Sprintf("%.2f", elapsedTime)
			zap.S().Infof("Fetching data... Elapsed Time: %s seconds", elapsedTimeFormatted)
		case err := <-errChan:
			return &pb.ResponseWeatherData{
				Message: "Error while fetching data",
				Success: false,
			}, err
		case weather := <-weatherChan:
			zap.S().Info("Data successfully fetched!")
			return weather, nil
		}
	}
}

func (s *WeatherServer) GetWeatherDataStream(req *pb.RequestWeatherData, stream pb.WeatherService_GetWeatherDataStreamServer) error {
	var wg sync.WaitGroup
	startTime := time.Now()
	ticker := time.NewTicker(2 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	weatherChan := make(chan *pb.ResponseStreamWeatherData)
	errChan := make(chan error)

	for {
		select {
		case <-ctx.Done():
			return errors.New("operation cancelled due to timeout")
		case <-ticker.C:
			zap.S().Info("Ticker ticked!")
			elapsedTime := time.Since(startTime).Seconds()
			elapsedTimeFormatted := fmt.Sprintf("%.2f", elapsedTime)

			wg.Add(1)
			go func() {
				defer wg.Done()

				city, err := s.Services.WeatherService.GetCityDataByUsername(req.Username)
				if city == "" || err != nil {
					city = *DEFAULT_CITY
				}

				city = strings.TrimSpace(city)
				time.Sleep(2 * time.Second)
				weatherData, err := s.Services.WeatherService.FetchWeatherData(city)
				if err != nil {
					errChan <- errors.New("weather data is empty")
					return
				}

				weather := &pb.Weather{
					City:        weatherData.City,
					Country:     weatherData.Country,
					Temperature: weatherData.TempC,
					CityTime:    weatherData.CityTime,
				}

				if weatherData.City == "" {
					errChan <- errors.New("weather data is empty")
					return
				}

				msg := fmt.Sprintf("Successfully fetched! Weather Temp: %v", weatherData.TempC)
				weatherChan <- &pb.ResponseStreamWeatherData{
					Weather: weather,
					Message: msg,
					Success: true,
					Time:    elapsedTimeFormatted,
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				data := <-weatherChan
				if data == nil {
					return
				}
				if err := stream.Send(data); err != nil {
					errChan <- err
				}
				zap.S().Info("Data successfully sended!")
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				wg.Wait()
				close(weatherChan)
				close(errChan)
			}()

		case err := <-errChan:
			return err
		case <-ctx.Done():
			return errors.New("operation finished due to timeout")
		}
	}
}
