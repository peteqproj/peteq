package main

import (
	"context"

	_ "github.com/lib/pq"

	"github.com/peteqproj/peteq/pkg/db/postgres"
	eventbus "github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
)

func main() {
	logr := logger.New(logger.Options{})

	db, err := postgres.Connect(utils.GetEnvOrDie("POSTGRES_URL"))
	utils.DieOnError(err, "Failed to connect to postgres")
	defer db.Close()

	ebus, err := eventbus.New(eventbus.Options{
		Type:       "rabbitmq",
		Logger:     logr.Fork("module", "eventbus"),
		EventlogDB: db,
		RabbitMQ: eventbus.RabbitMQOptions{
			Host:     utils.GetEnvOrDie("RABBITMQ_HOST"),
			Port:     utils.GetEnvOrDie("RABBITMQ_PORT"),
			APIPort:  utils.GetEnvOrDie("RABBITMQ_API_PORT"),
			Username: utils.GetEnvOrDie("RABBITMQ_USERNAME"),
			Password: utils.GetEnvOrDie("RABBITMQ_PASSWORD"),
		},
	})
	utils.DieOnError(err, "Failed to connect to eventbus")
	defer ebus.Stop()
	ctx := context.WithValue(context.Background(), "UserID", "e0a8a226-e64b-4e08-ac38-298771fdf71d")
	err = ebus.Replay(ctx)
	utils.DieOnError(err, "Failed to replay events")
}
