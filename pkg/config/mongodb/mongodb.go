package mongodb

import (
	"context"
	config_builder "github.com/bulutcan99/grpc_weather/pkg/config"
	"github.com/bulutcan99/grpc_weather/pkg/env"
	"github.com/bulutcan99/grpc_weather/pkg/utility"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	doOnce   sync.Once
	client   *mongo.Client
	database *mongo.Database
	DB_NAME  = &env.Env.DbName
)

type Mongo struct {
	client   *mongo.Client
	context  context.Context
	database *mongo.Database
}

func NewConnetion() *Mongo {
	ctx := context.Background()
	mongoCon, err := config_builder.ConnectionURLBuilder("mongo")
	if err != nil {
		panic(err)
	}

	doOnce.Do(func() {
		cli, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoCon))
		if err != nil {
			panic(err)
		}

		err = cli.Ping(ctx, nil)
		if err != nil {
			panic(err)
		}

		database = cli.Database(*DB_NAME)
		client = cli
	})

	zap.S().Debug("Connected to MongoDB successfully: ", mongoCon)

	return &Mongo{
		client:   client,
		context:  ctx,
		database: database,
	}
}

func (m *Mongo) Close() {
	err := m.client.Disconnect(m.context)
	if err != nil {
		zap.S().Errorf("Error while disconnecting from MongoDB: %s", err)
	}
	zap.S().Debug("Connection to MongoDB closed successfully")
}

func (m *Mongo) Stop() {
	mongodbActiveSessionsCount := m.numberSessionsInProgress()
	for mongodbActiveSessionsCount != 0 {
		mongodbActiveSessionsCount = m.numberSessionsInProgress()
	}
}

func (m *Mongo) numberSessionsInProgress() int {
	return m.client.NumberSessionsInProgress()
}

func (m *Mongo) CreateCollectionIfNotExist(collectionName string) *mongo.Collection {
	collectionList, err := m.database.ListCollectionNames(m.context, bson.D{{}})
	if err != nil {
		panic(err)
	}

	isCollectionExist, _ := utility.Contains(collectionList, collectionName, nil)
	if isCollectionExist {
		return m.database.Collection(collectionName)
	}

	err = m.database.CreateCollection(m.context, collectionName)
	if err != nil {
		panic(err)
	}

	return m.database.Collection(collectionName)
}

func (m *Mongo) GetContext() context.Context {
	return m.context
}
