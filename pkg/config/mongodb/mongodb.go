package config_mongodb

import (
	"context"
	config_builder "github.com/bulutcan99/grpc_weather/pkg/config"
	"github.com/bulutcan99/grpc_weather/pkg/env"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	doOnce  sync.Once
	client  *mongo.Client
	DB_NAME = &env.Env.DbName
)

type Mongo struct {
	Client   *mongo.Client
	Context  context.Context
	Database string
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
		client = cli
	})

	zap.S().Debug("Connected to MongoDB successfully: ", mongoCon)
	return &Mongo{
		Client:   client,
		Context:  ctx,
		Database: *DB_NAME,
	}
}

func (m *Mongo) Close() {
	err := m.Client.Disconnect(m.Context)
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
	return m.Client.NumberSessionsInProgress()
}
