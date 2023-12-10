package service

import (
	"errors"
	"github.com/bulutcan99/grpc_weather/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IUserRepo interface {
	Insert(user model.User) (primitive.ObjectID, error)
	FindOne(filter any, opts *options.FindOneOptions) (*model.User, error)
	UpdateOne(filter any, update any) (*model.User, error)
	DeleteOne(filter any) (*mongo.DeleteResult, error)
}

type UserService struct {
	UserRepo IUserRepo
}

func NewUserService(userRepo IUserRepo) *UserService {
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

	resultID, err := u.UserRepo.Insert(newUser)
	if err != nil {
		return primitive.NilObjectID, errors.New("user is not registered")
	}

	return resultID, nil
}

func (u *UserService) LoginUser(username string, password string) (*model.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username or password is empty")
	}

	filter := bson.D{
		{"username", username},
		{"password", password},
	}

	opts := options.FindOne().SetProjection(bson.M{
		"user_id": 1,
	})
	user, err := u.UserRepo.FindOne(filter, opts)
	if err != nil {
		return nil, errors.New("user is not found")
	}

	if user == nil {
		return nil, errors.New("user is not found")
	}

	return user, nil
}

func (u *UserService) UpdateUserPassword(id primitive.ObjectID, pass string) (*model.User, error) {
	filter := bson.D{
		{"_id", id},
	}

	update := bson.D{
		{"$set", bson.D{
			{"password", pass},
		}},
	}

	user, err := u.UserRepo.UpdateOne(filter, update)
	if err != nil {
		return nil, errors.New("user is not updated")
	}

	return user, nil
}

func (u *UserService) DeleteUser(id primitive.ObjectID) error {
	if id == primitive.NilObjectID {
		return errors.New("user id is empty")
	}

	filter := bson.D{
		{"_id", id},
	}

	_, err := u.UserRepo.DeleteOne(filter)
	if err != nil {
		return errors.New("user is not deleted")
	}

	return nil
}
