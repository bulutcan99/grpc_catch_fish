package grpc_server

import (
	"errors"
	"github.com/bulutcan99/grpc_weather/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(user model.User) (primitive.ObjectID, error) {
	args := m.Called(user)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *MockUserService) FindUser(username string, password string) (*model.User, error) {
	args := m.Called(username, password)
	return args.Get(0).(*model.User), args.Error(1)
}

func TestRegisterUser(t *testing.T) {
	mockUserService := new(MockUserService)
	weatherServer := NewWeatherServer(mockUserService)
	testCases := []struct {
		name          string
		user          model.User
		expectedID    primitive.ObjectID
		expectedError error
	}{
		{
			name: "ValidUser",
			user: model.User{
				Username: "john_doe",
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
				City:     "New York",
			},
			expectedID:    primitive.NewObjectID(),
			expectedError: nil,
		},
		{
			name:          "EmptyUser",
			user:          model.User{},
			expectedID:    primitive.NilObjectID,
			expectedError: errors.New("user is empty"),
		},
		{
			name: "UserWithMissingFields",
			user: model.User{
				Username: "john_doe",
				Name:     "John Doe",
				Email:    "",
				Password: "password123",
				City:     "New York",
			},
			expectedID:    primitive.NilObjectID,
			expectedError: errors.New("user with missing fields"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserService.On("RegisterUser", tc.user).Return(tc.expectedID, tc.expectedError)
			resultID, err := weatherServer.UserService.RegisterUser(tc.user)
			if tc.expectedError == nil {
				require.NoError(t, err)
				require.NotNil(t, resultID)
			} else {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
			}
			mockUserService.AssertExpectations(t)
		})
	}
}
