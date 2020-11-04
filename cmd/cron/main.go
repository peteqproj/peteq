package main

import (
	triggerDomain "github.com/peteqproj/peteq/domain/trigger"
	userDomain "github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/pkg/config"
	"github.com/peteqproj/peteq/pkg/cron"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/db/postgres"
	eventbus "github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/server"
	"github.com/peteqproj/peteq/pkg/utils"
)

func main() {
	logr := logger.New(logger.Options{})
	cnf := &config.Server{
		Port: utils.GetEnvOrDie("PORT"),
	}
	s := server.New(server.Options{
		Config: cnf,
	})

	pg, err := postgres.Connect(utils.GetEnvOrDie("POSTGRES_URL"))
	defer pg.Close()
	db := db.New(db.Options{
		DB: pg,
	})
	utils.DieOnError(err, "Failed to connect to postgres")

	userRepo := &userDomain.Repo{
		DB:     db,
		Logger: logr.Fork("repo", "user"),
	}
	users, err := userRepo.List(userDomain.ListOptions{})
	utils.DieOnError(err, "Failed to load users")

	triggerRepo := &triggerDomain.Repo{
		DB:     db,
		Logger: logr.Fork("repo", "trigger"),
	}

	ebus, err := eventbus.New(eventbus.Options{
		Type:        "rabbitmq",
		Logger:      logr.Fork("module", "eventbus"),
		EventlogDB:  db,
		WatchQueues: false,
		RabbitMQ: eventbus.RabbitMQOptions{
			Host:     utils.GetEnvOrDie("RABBITMQ_HOST"),
			Port:     utils.GetEnvOrDie("RABBITMQ_PORT"),
			APIPort:  utils.GetEnvOrDie("RABBITMQ_API_PORT"),
			Username: utils.GetEnvOrDie("RABBITMQ_USERNAME"),
			Password: utils.GetEnvOrDie("RABBITMQ_PASSWORD"),
		},
	})
	utils.DieOnError(err, "Failed to connect to eventbus")
	err = ebus.Start()
	utils.DieOnError(err, "Failed to start eventbus")
	defer ebus.Stop()

	logr.Info("All user loaded", "len", len(users))
	for _, user := range users {
		triggers, err := triggerRepo.List(triggerDomain.QueryOptions{
			UserID: user.Metadata.ID,
		})
		utils.DieOnError(err, "Failed to load users triggers")
		logr.Info("All user triggers loaded", "len", len(triggers))
		cr := cron.New(cron.Options{
			EventPublisher: ebus,
			Logger:         logr.Fork("user", user.Metadata.ID),
			UserID:         user.Metadata.ID,
		})
		for _, t := range triggers {
			if t.Spec.Cron != nil {
				logr.Info("Starting to watch", "user", user.Metadata.Email, "cron", t.Spec.Cron)
				cr.AddFunc(t.Metadata.ID, *t.Spec.Cron)
			}
		}
		go cr.Start()
		defer cr.Stop()
	}

	err = s.Start()
	utils.DieOnError(err, "Failed to run server")
}
