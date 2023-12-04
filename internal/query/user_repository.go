package query

import (
	"context"
	"errors"
	"github.com/bulutcan99/grpc_weather/model"
	config_mongodb "github.com/bulutcan99/grpc_weather/pkg/config/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
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
	ID := primitive.NewObjectIDFromTimestamp(time.Now())
	newUser := model.User{
		ID:       ID,
		Username: user.Username,
		Password: user.Password,
		Name:     user.Name,
		Email:    user.Email,
		City:     user.City,
	}

	doc, err := u.userCollection.InsertOne(u.ctx, newUser)
	if err != nil {
		return primitive.NilObjectID, errors.New("user is not registered")
	}

	return doc.InsertedID.(primitive.ObjectID), nil
}

func (u *UserRepository) FindOne(username string, password string) (*model.User, error) {
	var user model.User
	filter := bson.D{
		{Key: "username", Value: username},
		{Key: "password", Value: password},
	}
	data := u.userCollection.FindOne(u.ctx, filter)
	err := data.Decode(&user)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("user is not found")
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
