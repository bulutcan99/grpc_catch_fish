package grpc_client

import (
	"context"
	"github.com/bulutcan99/grpc_weather/internal/fetch"
	"github.com/bulutcan99/grpc_weather/proto/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"math/rand"
	"sync"
	"time"
)

type User struct {
	Lat  float64
	Long float64
	City string
}

type WeatherClient struct {
	*User
	conn    *grpc.ClientConn
	client  pb.WeatherServiceClient
	Fetcher *fetch.FetchingDataClient
	mutex   *sync.Mutex
}

func NewWeatherClient() *WeatherClient {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Fatal(err)
	}

	return &WeatherClient{
		mutex:   new(sync.Mutex),
		conn:    conn,
		Fetcher: fetch.NewFetchingDataClient(),
	}
}

func NewUser() *User {
	return &User{
		Lat:  randFloat(-90, 90),
		Long: randFloat(-180, 180),
	}
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func (c *WeatherClient) GetCity() (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	url := c.Fetcher.GetCityUrl(c.User.Lat, c.User.Long)
	city, err := c.Fetcher.FetchCity(url)
	if err != nil {
		return "", err
	}
	return city, nil
}

func (c *WeatherClient) Close() error {
	return c.conn.Close()
}

func (c *WeatherClient) GetWeatherDataByLatLong() error {
	ctx := context.Background()
	stream, err := c.client.GetWeatherDataByLatLongStream(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error)
	wg.Add(2)
	go c.City(stream, errChan, &wg)
	go c.WeatherData(stream, errChan, &wg)
	wg.Wait()
	stream.CloseSend()
	for {
		select {
		case err := <-errChan:
			if err != nil {
				zap.S().Error(err)
			}
			return err
		default:
			return nil
		}
	}
}

func (c *WeatherClient) WeatherData(stream pb.WeatherService_GetWeatherDataByLatLongStreamClient, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-stream.Context().Done():
			errChan <- nil
			return

		default:
			in, err := stream.Recv()
			if err == io.EOF {
				errChan <- nil
				return
			}
			zap.S().Info("Temperature: ", in.Temperature)

			if in.Status {
				zap.S().Info("Finished")
				stream.Context().Done()
				errChan <- err
				return
			}
		}
	}
}

func (c *WeatherClient) City(stream pb.WeatherService_GetWeatherDataByLatLongStreamClient, errChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	userCh := make(chan *User)
	go func() {
		for {
			select {
			case <-ticker.C:
				user := NewUser()
				c.User = user
				city, err := c.GetCity()
				if err != nil {
					zap.S().Error(err)
					errChan <- err
				}
				c.User.City = city
				userCh <- user
			}
		}
	}()

	for {
		select {
		case user := <-userCh:
			zap.S().Info("Getting city...")
			err := stream.Send(&pb.RequestStreamUserByLatLong{
				City: user.City,
			})
			if err != nil {
				zap.S().Error("Error while sending data: ", err)
				errChan <- err
			}
		}
	}
}

func (c *WeatherClient) GetWeatherDataStream(ctx context.Context, in *pb.RequestWeatherData, opts ...grpc.CallOption) (pb.WeatherService_GetWeatherDataStreamClient, error) {
	return c.client.GetWeatherDataStream(ctx, in, opts...)
}

func (c *WeatherClient) GetWeatherData(ctx context.Context, in *pb.RequestWeatherData, opts ...grpc.CallOption) (*pb.ResponseWeatherData, error) {
	return c.client.GetWeatherData(ctx, in, opts...)
}

func (c *WeatherClient) GetUserCity(ctx context.Context, in *pb.RequstUserCity, opts ...grpc.CallOption) (*pb.ResponseUserCity, error) {
	return c.client.GetUserCity(ctx, in, opts...)
}
