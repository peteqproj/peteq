package cmd

import (
	"context"
	"time"

	userDomain "github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/internal"

	sensorDomain "github.com/peteqproj/peteq/domain/sensor"
	"github.com/peteqproj/peteq/pkg/config"
	"github.com/peteqproj/peteq/pkg/cron"
	"github.com/peteqproj/peteq/pkg/db"
	eventbus "github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/server"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"

	"github.com/spf13/cobra"
)

var cronServiceCmdFlags struct {
	verbose          bool
	postgresURL      string
	rabbitmqHost     string
	rabbitmqPort     string
	rabbitmqAPIPort  string
	rabbitmqUsername string
	rabbitmqPassword string
	port             string
}

type (
	userSensorPair struct {
		user    userDomain.User
		sensors map[string]*sensorDomain.Sensor
		cron    cron.Cron
	}
)

var cronServiceCmd = &cobra.Command{
	Use:   "cron-service",
	Short: "Starts the cron server",
	RunE: func(cmd *cobra.Command, args []string) error {
		logr := logger.New(logger.Options{})
		cnf := &config.Server{
			Port: utils.GetEnvOrDie("PORT"),
		}
		s := server.New(server.Options{
			Config: cnf,
		})

		db, err := db.New(db.Options{
			URL: utils.GetEnvOrDie("POSTGRES_URL"),
		})
		utils.DieOnError(err, "Failed to connect to postgres")
		sqldb, err := db.DB()
		utils.DieOnError(err, "Failed to get sqldb object from gorm")
		defer sqldb.Close()

		userRepo := &userDomain.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "user"),
		}

		sensorRepo := &sensorDomain.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "sensor"),
		}
		if err := sensorRepo.Initiate(context.Background()); err != nil {
			utils.DieOnError(err, "Failed to init sensor repo")
		}
		ebus := internal.NewEventBusFromFlagsOrDie(db, userRepo, false, logr.Fork("module", "eventbus"))
		if err := ebus.Start(); err != nil {
			return err
		}
		defer ebus.Stop()
		logr.Info("Eventbus connected")

		go loop(userRepo, sensorRepo, ebus, logr)

		s.SetReady()
		err = s.Start()
		return err
	},
}

func init() {
	startCmd.AddCommand(cronServiceCmd)
}
func loop(userRepo *userDomain.Repo, sensorRepo *sensorDomain.Repo, ebus eventbus.Eventbus, lgr logger.Logger) {
	l := map[string]userSensorPair{}
	for {
		select {
		case _ = <-time.After(time.Minute * 1):
			{
				lgr.Info("Running loop")
				res, err := userRepo.List(context.Background())
				if err != nil {
					lgr.Info("Failed to load users", "error", err.Error())
					continue
				}
				for _, u := range res {
					if _, found := l[u.Metadata.ID]; found {
						continue
					}
					lgr.Info("New user added", "email", u.Spec.Email, "id", u.Metadata.ID)
					l[u.Metadata.ID] = userSensorPair{
						user:    *u,
						sensors: map[string]*sensorDomain.Sensor{},
						cron: cron.New(cron.Options{
							EventPublisher: ebus,
							Logger:         lgr.Fork("user", u.Metadata.ID),
							UserID:         u.Metadata.ID,
						}),
					}
					go l[u.Metadata.ID].cron.Start()
				}

				for id, pair := range l {
					ctx := tenant.ContextWithUser(context.Background(), pair.user)
					res, err := sensorRepo.ListByUserid(ctx, id)
					if err != nil {
						lgr.Info("Failed to load user sensors", "error", err.Error(), "user", id)
					}

					for _, t := range res {
						if _, found := l[id].sensors[t.Metadata.ID]; found {
							continue
						}
						if t.Spec.Cron == nil {
							continue
						}
						lgr.Info("New cron sensor added", "cron", *t.Spec.Cron, "id", t.Metadata.ID)
						pair.sensors[t.Metadata.ID] = t
						pair.cron.AddFunc(t.Metadata.ID, *t.Spec.Cron)
					}
				}
			}
		}
	}
}
