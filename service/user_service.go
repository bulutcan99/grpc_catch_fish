package service

import (
	"errors"
	"github.com/bulutcan99/grpc_weather/internal/query"
	"github.com/bulutcan99/grpc_weather/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUserService interface {
	RegisterUser(newUser model.User) (primitive.ObjectID, error)
}

type UserService struct {
	UserRepo query.IUserRepo
}

func NewUserService(userRepo query.IUserRepo) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (u *UserService) RegisterUser(newUser model.User) (primitive.ObjectID, error) {
	if newUser == (model.User{}) {
		return primitive.NilObjectID, errors.New("user is empty")
	}
	if newUser.Username == "" || newUser.Name == "" || newUser.Email == "" || newUser.Password == "" || newUser.City == "" {
		return primitive.NilObjectID, errors.New("user with missing fields")
	}

	resultID, err := u.UserRepo.Register(newUser)
	if err != nil {
		return primitive.NilObjectID, errors.New("user is not registered")
	}

	return resultID, nil
}

func (u *UserService) FindUser(username string, password string) (*model.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username or password is empty")
	}

	user, err := u.UserRepo.Find(username, password)
	if err != nil {
		return nil, errors.New("user is not found")
	}

	return &user, nil
}
