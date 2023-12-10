package query

import (
	"context"
	config_mongodb "github.com/bulutcan99/grpc_weather/pkg/config/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type WeatherRepository struct {
	client            *mongo.Client
	ctx               context.Context
	weatherCollection *mongo.Collection
}

func NewWeatherRepository(mongo *config_mongodb.Mongo, collectionName string) *WeatherRepository {
	weatherCollection := mongo.Client.Database(mongo.Database).Collection(collectionName)
	return &WeatherRepository{
		client:            mongo.Client,
		ctx:               mongo.Context,
		weatherCollection: weatherCollection,
	}
}

func (w *WeatherRepository) InsertOrUpdate(filter any, update any) *mongo.SingleResult {
	update = bson.M{
		"$set": update,
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}

	upsertOption := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	response := w.weatherCollection.FindOneAndUpdate(w.ctx, filter, update, upsertOption)

	return response

}

func (w *WeatherRepository) FindOne(filter any) *mongo.SingleResult {
	response := w.weatherCollection.FindOne(w.ctx, filter)
	return response
}
