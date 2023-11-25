package service

import (
	"github.com/bulutcan99/grpc_weather/internal/query"
	"github.com/bulutcan99/grpc_weather/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	UserRepo query.UserRepo
}

func NewUserService(userRepo query.UserRepo) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (u *UserService) RegisterUser(user model.User) (primitive.ObjectID, error) {
	if user.Username == "" || user.Name == "" || user.Email == "" || user.Password == "" || user.City == "" {
		return primitive.ObjectID{}, nil
	}

	return u.UserRepo.RegisterUser(user)
}

func (u *UserService) FindUser(username string, password string) (model.User, error) {
	if username == "" || password == "" {
		return model.User{}, nil
	}

	return u.UserRepo.FindUser(username, password)
}
