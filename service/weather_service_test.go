package service

import (
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestWeatherService_FetchWeatherData(t *testing.T) {
	testCases := []struct {
		name        string
		inputCity   string
		expectedErr error
	}{
		{
			name:        "FetchWeatherData_Success",
			inputCity:   "Istanbul",
			expectedErr: nil,
		},
		{
			name:        "FetchWeatherData_EmptyCity",
			inputCity:   "",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			weatherRepo := &mockWeatherRepo{}
			userRepo := &mockUserRepo{}
			weatherService := NewWeatherService(weatherRepo, userRepo)

			_, err := weatherService.FetchWeatherData(tc.inputCity)

			if err != nil {
				t.Errorf("error while fetching weather data: %v", err)
			}
		})
	}
}

type mockWeatherRepo struct {
}

func (m *mockWeatherRepo) InsertOrUpdate(filter interface{}, update interface{}) *mongo.SingleResult {
	return nil
}

func (m *mockWeatherRepo) FindOne(filter interface{}) *mongo.SingleResult {
	return nil
}
