package service

import (
	"errors"
	"github.com/bulutcan99/grpc_weather/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Register(user model.User) (primitive.ObjectID, error) {
	args := m.Called(user)
	return args.Get(0).(primitive.ObjectID), args.Error(1)
}

func (m *MockUserRepo) Find(username string, password string) (model.User, error) {
	args := m.Called(username, password)
	return args.Get(0).(model.User), args.Error(1)

}

func TestRegister(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	userService := NewUserService(mockUserRepo)
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
			mockUserRepo.On("Register", tc.user).Return(tc.expectedID, tc.expectedError)
			resultID, err := userService.RegisterUser(tc.user)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedID, resultID)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestFindUser(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	userService := NewUserService(mockUserRepo)
	testCases := []struct {
		name          string
		username      string
		password      string
		expectedUser  model.User
		expectedError error
	}{
		{
			name:     "ValidUser",
			username: "john_doe",
			password: "password123",
			expectedUser: model.User{
				Username: "john_doe",
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
				City:     "New York",
			},
			expectedError: nil,
		},
		{
			name:          "UserWithEmptyFields",
			username:      "",
			password:      "invalid_password",
			expectedUser:  model.User{},
			expectedError: errors.New("username or password is empty"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo.On("FindUser", tc.username, tc.password).Return(tc.expectedUser, tc.expectedError)
			resultUser, err := userService.FindUser(tc.username, tc.password)
			if tc.expectedError == nil {
				require.NoError(t, err)
				require.NotNil(t, resultUser)
				require.NotEmpty(t, resultUser.City)
				if len(tc.expectedUser.City) > 0 {
					require.Equal(t, tc.expectedUser.City, resultUser.City)
				}
			} else {
				require.Error(t, err)
				require.Nil(t, resultUser)
				require.Equal(t, tc.expectedError, err)
			}
		})
	}
}
