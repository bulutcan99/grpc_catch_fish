package query

import (
	"testing"

	"github.com/bulutcan99/grpc_weather/model"
	"github.com/bulutcan99/grpc_weather/pkg/config/mongodb"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserRepository_Insert(t *testing.T) {
	mongoConfig := config_mongodb.NewConnetion()

	userRepo := NewUserRepositry(mongoConfig, "users")

	testCases := []struct {
		name        string
		inputUser   model.User
		expectedID  primitive.ObjectID
		expectedErr error
	}{
		{
			name: "Insert_Success",
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
			name: "Insert_DuplicateUsername",
			inputUser: model.User{
				Username: "existinguser",
				Name:     "Existing User",
				Email:    "existing@example.com",
				Password: "existingpassword",
				City:     "ExistingCity",
			},
			expectedID:  primitive.NilObjectID,
			expectedErr: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultID, err := userRepo.Insert(tc.inputUser)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedID, resultID)
		})
	}
}

func TestUserRepository_FindOne(t *testing.T) {
	mongoConfig := config_mongodb.NewConnetion()

	userRepo := NewUserRepositry(mongoConfig, "users")

	testCases := []struct {
		name         string
		username     string
		password     string
		expectedUser *model.User
		expectedErr  error
	}{
		{
			name:     "FindOne_UserExists",
			username: "existinguser",
			password: "existingpassword",
			expectedUser: &model.User{
				Username: "existinguser",
				Name:     "Existing User",
				Email:    "existing@example.com",
				Password: "existingpassword",
				City:     "ExistingCity",
			},
			expectedErr: nil,
		},
		{
			name:         "FindOne_UserNotFound",
			username:     "nonexistentuser",
			password:     "password",
			expectedUser: nil,
			expectedErr:  assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			foundUser, err := userRepo.FindOne(tc.username, tc.password)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedUser, foundUser)
		})
	}
}
