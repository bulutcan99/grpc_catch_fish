package query

import (
	"context"
	"github.com/bulutcan99/grpc_weather/model"
	config_mongodb "github.com/bulutcan99/grpc_weather/pkg/config/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserRepo interface {
	Insert(user model.User) (primitive.ObjectID, error)
	FindOne(username string, password string) (*model.User, error)
}

type UserRepository struct {
	client         *mongo.Client
	ctx            context.Context
	userCollection *mongo.Collection
}

func NewUserRepositry(mongo *config_mongodb.Mongo, collectionName string) *UserRepository {
	userCollection := mongo.Client.Database(mongo.Database).Collection(collectionName)
	return &UserRepository{
		client:         mongo.Client,
		ctx:            mongo.Context,
		userCollection: userCollection,
	}
}

func (u *UserRepository) Insert(user model.User) (primitive.ObjectID, error) {
	res, err := u.userCollection.InsertOne(u.ctx, user)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (u *UserRepository) FindOne(username string, password string) (*model.User, error) {
	var user model.User
	err := u.userCollection.FindOne(u.ctx, model.User{Username: username, Password: password}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
