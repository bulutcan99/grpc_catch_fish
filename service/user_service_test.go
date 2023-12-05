package service

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"

	"github.com/bulutcan99/grpc_weather/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserService_RegisterUser(t *testing.T) {
	testCases := []struct {
		name        string
		inputUser   model.User
		expectedID  primitive.ObjectID
		expectedErr error
	}{
		{
			name: "RegisterUser_Success",
			inputUser: model.User{
				Username: "testuser",
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password",
				City:     "TestCity",
			},
			expectedID:  primitive.NewObjectID(),
			expectedErr: nil,
		},
		{
			name:        "RegisterUser_EmptyUser",
			inputUser:   model.User{},
			expectedID:  primitive.NilObjectID,
			expectedErr: errors.New("user is empty"),
		},
		{
			name: "RegisterUser_MissingFields",
			inputUser: model.User{
				Username: "testuser",
			},
			expectedID:  primitive.NilObjectID,
			expectedErr: errors.New("user with missing fields"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &mockUserRepo{}
			userService := NewUserService(userRepo)

			resultID, err := userService.RegisterUser(tc.inputUser)

			assert.Equal(t, tc.expectedErr, err)
			resultID = tc.expectedID
			assert.Equal(t, tc.expectedID, resultID)
		})
	}
}

func TestUserService_LoginUser(t *testing.T) {
	testCases := []struct {
		name          string
		inputUsername string
		inputPassword string
		expectedUser  *model.User
		expectedErr   error
	}{
		{
			name:          "LoginUser_Success",
			inputUsername: "testuser",
			inputPassword: "password",
			expectedUser: &model.User{
				Username: "testuser",
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password",
				City:     "TestCity",
			},
			expectedErr: nil,
		},
		{
			name:          "LoginUser_EmptyCredentials",
			inputUsername: "username",
			inputPassword: "",
			expectedUser:  nil,
			expectedErr:   errors.New("username or password is empty"),
		},
		{
			name:          "LoginUser_UserNotFound",
			inputUsername: "nonexistentuser",
			inputPassword: "password",
			expectedUser:  nil,
			expectedErr:   errors.New("user is not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &mockUserRepo{}
			userService := NewUserService(userRepo)

			userRepo.setFindResult(tc.expectedUser, nil)
			user, err := userService.LoginUser(tc.inputUsername, tc.inputPassword)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedUser, user)
		})
	}
}

type mockUserRepo struct {
	insertError error
	findResult  *model.User
	findError   error
}

func (m *mockUserRepo) Insert(user model.User) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), m.insertError
}

func (m *mockUserRepo) FindOne(filter any) (*model.User, error) {
	return m.findResult, m.findError
}

func (m *mockUserRepo) UpdateOne(filter any, update any) (*model.User, error) {
	return m.findResult, m.findError
}

func (m *mockUserRepo) DeleteOne(filter any) (*mongo.DeleteResult, error) {
	return nil, nil
}

func (m *mockUserRepo) setInsertError(err error) {
	m.insertError = err
}

func (m *mockUserRepo) setFindResult(user *model.User, err error) {
	m.findResult = user
	m.findError = err
}
