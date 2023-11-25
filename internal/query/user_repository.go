package query

import (
	"context"
	"fmt"
	"github.com/bulutcan99/grpc_weather/model"
	config_mongodb "github.com/bulutcan99/grpc_weather/pkg/config/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo interface {
	RegisterUser(user model.User) (primitive.ObjectID, error)
	FindUser(username string, password string) (model.User, error)
}

type UserRepositry struct {
	client         *mongo.Client
	ctx            context.Context
	userCollection *mongo.Collection
}

func NewUserRepositry(mongo *config_mongodb.Mongo, collectionName string) *UserRepositry {
	userCollection := mongo.Client.Database(mongo.Database).Collection(collectionName)
	fmt.Println("userCollection: ", userCollection)
	return &UserRepositry{
		client:         mongo.Client,
		ctx:            mongo.Context,
		userCollection: userCollection,
	}
}

func (u *UserRepositry) RegisterUser(user model.User) (primitive.ObjectID, error) {
	res, err := u.userCollection.InsertOne(u.ctx, user)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (u *UserRepositry) FindUser(username string, password string) (model.User, error) {
	var user model.User
	err := u.userCollection.FindOne(u.ctx, model.User{Username: username, Password: password}).Decode(&user)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
