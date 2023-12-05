package query

import (
	"context"
	"errors"
	"github.com/bulutcan99/grpc_weather/model"
	config_mongodb "github.com/bulutcan99/grpc_weather/pkg/config/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type IUserRepo interface {
	Insert(user model.User) (primitive.ObjectID, error)
	FindOne(filter any) (*model.User, error)
	UpdateOne(filter any, update any) (*model.User, error)
	DeleteOne(filter any) (*mongo.DeleteResult, error)
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

func (u *UserRepository) FindOne(filter any) (*model.User, error) {
	var user model.User
	response := u.userCollection.FindOne(u.ctx, filter)
	if err := response.Decode(&user); err != nil {
		return nil, errors.New("user is not found")
	}

	return &user, nil
}

func (u *UserRepository) Find(filter any) (*mongo.Cursor, error) {
	response, err := u.userCollection.Find(u.ctx, filter)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (u *UserRepository) UpdateOne(filter any, update any) (*model.User, error) {
	var user model.User
	res := u.userCollection.FindOneAndUpdate(u.ctx, filter, update)
	if err := res.Decode(&user); err != nil {
		return nil, errors.New("user is not updated")
	}

	return &user, nil
}

func (u *UserRepository) Update(user model.User, update any) (*mongo.UpdateResult, error) {
	res, err := u.userCollection.UpdateOne(u.ctx, user, update)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *UserRepository) DeleteOne(filter any) (*mongo.DeleteResult, error) {
	res, err := u.userCollection.DeleteOne(u.ctx, filter)
	if err != nil {
		return nil, err
	}
	return res, nil
}
