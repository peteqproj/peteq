package internal

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/peteqproj/peteq/domain/user"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/event"
	eventbus "github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/event/storage"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
	"google.golang.org/api/option"
)

func NewEventBusFromFlagsOrDie(db db.Database, repo *user.Repo, watchQueues bool, logger logger.Logger) eventbus.Eventbus {
	logger.Info("Connecting to eventbus")
	etype := utils.GetEnvOrDie("EVENTBUS_TYPE")
	var rabbit *eventbus.RabbitMQOptions
	var google *eventbus.GooglePubSubOptions
	if etype == "rabbitmq" {
		rabbit = &eventbus.RabbitMQOptions{
			WatchQueues: watchQueues,
			Host:        utils.GetEnvOrDie("RABBITMQ_HOST"),
			Port:        utils.GetEnvOrDie("RABBITMQ_PORT"),
			APIPort:     utils.GetEnvOrDie("RABBITMQ_API_PORT"),
			Username:    utils.GetEnvOrDie("RABBITMQ_USERNAME"),
			Password:    utils.GetEnvOrDie("RABBITMQ_PASSWORD"),
		}
	}

	if etype == "google" {
		c, err := pubsub.NewClient(context.Background(), "peteq-291604", option.WithCredentialsFile(utils.GetEnvOrDie("GOOGLE_PUBSUB_EVENT_BUS_SA_CREDENTIALS")))
		utils.DieOnError(err, "Failed to create Google pub-sub client")
		google = &eventbus.GooglePubSubOptions{
			Client: c,
		}
	}
	bus, err := eventbus.New(eventbus.Options{
		Logger: logger,
		ExtendContextFunc: func(ctx context.Context, ev event.Event) context.Context {
			if ev.Tenant.Type != tenant.User.String() {
				return ctx
			}
			user, err := repo.GetById(ctx, ev.Tenant.ID)
			if err != nil {
				logger.Info("Failed extend context", "user", ev.Tenant.ID, "event", ev.Metadata.ID)
				return ctx
			}
			return tenant.ContextWithUser(ctx, *user)
		},
		EventStorage: storage.New(storage.Options{
			DB: db,
		}),
		RabbitMQ:     rabbit,
		GooglePubSub: google,
	})
	utils.DieOnError(err, "Failed to connect to eventbus")
	return bus
}

func NewCommandBusFromFlagsOrDie(repo *user.Repo, logger logger.Logger) commandbus.CommandBus {
	logger.Info("Connecting to commandbus")
	etype := utils.GetEnvOrDie("COMMANDBUS_TYPE")
	var google *commandbus.GoogleCommandBusOptions
	if etype == "google" {
		c, err := pubsub.NewClient(context.Background(), "peteq-291604", option.WithCredentialsFile(utils.GetEnvOrDie("GOOGLE_PUBSUB_COMMAND_BUS_SA_CREDENTIALS")))
		utils.DieOnError(err, "Failed to create Google pub-sub client")
		google = &commandbus.GoogleCommandBusOptions{
			PubSubClient:       c,
			PubSubTopic:        utils.GetEnvOrDie("GOOGLE_PUBSUB_COMMAND_BUS_TOPIC"),
			PubSubSubscribtion: utils.GetEnvOrDie("GOOGLE_PUBSUB_COMMAND_BUS_TOPIC_SUBSCRIBTION"),
		}
	}

	return commandbus.New(commandbus.Options{
		GoogleOptions: google,
		ExtendContextFunc: func(c context.Context, id string) context.Context {
			user, err := repo.GetById(c, id)
			if err != nil {
				logger.Info("Failed extend context", "user", id)
				return c
			}
			return tenant.ContextWithUser(c, *user)
		},
		Logger: logger,
	})
}
