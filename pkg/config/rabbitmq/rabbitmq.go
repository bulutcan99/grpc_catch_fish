package config_rabbitmq

import (
	"context"
	config_builder "github.com/bulutcan99/grpc_weather/pkg/config"
	"sync"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var (
	once           sync.Once
	rabbitmqClient *amqp.Connection
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Context    context.Context
}

func NewRabbitMQConnection() *RabbitMQ {
	ctx := context.Background()
	rabbitCon, err := config_builder.ConnectionURLBuilder("rabbitmq")
	if err != nil {
		panic(err)
	}

	once.Do(func() {
		conn, err := amqp.Dial(rabbitCon)
		if err != nil {
			panic(err)
		}

		rabbitmqClient = conn
	})

	zap.S().Info("Connected to RabbitMQ successfully.")
	return &RabbitMQ{
		Connection: rabbitmqClient,
		Context:    ctx,
	}
}

func (r *RabbitMQ) Close() {
	if err := r.Connection.Close(); err != nil {
		zap.S().Errorf("Error while closing the RabbitMQ connection: %s", err)
	}

	zap.S().Info("Connection to RabbitMQ closed successfully")
}

func (r *RabbitMQ) DeclareQueue(queueName string) error {
	ch, err := r.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
