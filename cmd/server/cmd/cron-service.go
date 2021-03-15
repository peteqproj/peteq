package cmd

import (
	"time"

	userDomain "github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/internal"

	triggerDomain "github.com/peteqproj/peteq/domain/trigger"
	"github.com/peteqproj/peteq/pkg/config"
	"github.com/peteqproj/peteq/pkg/cron"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/db/postgres"
	eventbus "github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/server"
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
	userTriggerPair struct {
		user     userDomain.User
		triggers map[string]triggerDomain.Trigger
		cron     cron.Cron
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

		triggerRepo := &triggerDomain.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "trigger"),
		}
		ebus := internal.NewEventBusFromFlagsOrDie(db, userRepo, false, logr.Fork("module", "eventbus"))
		if err := ebus.Start(); err != nil {
			return err
		}
		defer ebus.Stop()
		logr.Info("Eventbus connected")

		go loop(userRepo, triggerRepo, ebus, logr)

		s.SetReady()
		err = s.Start()
		return err
	},
}

func init() {
	startCmd.AddCommand(cronServiceCmd)
}
func loop(userRepo *userDomain.Repo, triggerRepo *triggerDomain.Repo, ebus eventbus.Eventbus, lgr logger.Logger) {
	l := map[string]userTriggerPair{}
	for {
		select {
		case _ = <-time.After(time.Minute * 1):
			{
				lgr.Info("Running loop")
				res, err := userRepo.List(userDomain.ListOptions{})
				if err != nil {
					lgr.Info("Failed to load users", "error", err.Error())
					continue
				}
				for _, u := range res {
					if _, found := l[u.Metadata.ID]; found {
						continue
					}
					lgr.Info("New user added", "email", u.Metadata.Email, "id", u.Metadata.ID)
					l[u.Metadata.ID] = userTriggerPair{
						user:     u,
						triggers: map[string]triggerDomain.Trigger{},
						cron: cron.New(cron.Options{
							EventPublisher: ebus,
							Logger:         lgr.Fork("user", u.Metadata.ID),
							UserID:         u.Metadata.ID,
						}),
					}
					go l[u.Metadata.ID].cron.Start()
				}

				for id, pair := range l {
					res, err := triggerRepo.List(triggerDomain.QueryOptions{
						UserID: id,
					})
					if err != nil {
						lgr.Info("Failed to load user triggers", "error", err.Error(), "user", id)
					}

					for _, t := range res {
						if _, found := l[id].triggers[t.Metadata.ID]; found {
							continue
						}
						if t.Spec.Cron == nil {
							continue
						}
						lgr.Info("New cron trigger added", "cron", *t.Spec.Cron, "id", t.Metadata.ID)
						pair.triggers[t.Metadata.ID] = t
						pair.cron.AddFunc(t.Metadata.ID, *t.Spec.Cron)
					}
				}
			}
		}
	}
}
