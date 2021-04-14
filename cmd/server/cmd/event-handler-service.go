package cmd

import (
	"context"

	automationDomain "github.com/peteqproj/peteq/domain/automation"
	listDomain "github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/task"
	triggerDomain "github.com/peteqproj/peteq/domain/trigger"
	userDomain "github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/internal"
	"github.com/peteqproj/peteq/pkg/config"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/db/postgres"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/repo"
	"github.com/peteqproj/peteq/pkg/server"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/peteqproj/peteq/saga"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var eventHandlerServiceCmdFlags struct {
	verbose          bool
	postgresURL      string
	rabbitmqHost     string
	rabbitmqPort     string
	rabbitmqAPIPort  string
	rabbitmqUsername string
	rabbitmqPassword string
	port             string
}

var eventHandlerServiceCmd = &cobra.Command{
	Use:   "event-handler-service",
	Short: "Starts the event handler service",
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
		taskRepo := task.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "task"),
		}
		if err := taskRepo.Initiate(context.Background()); err != nil {
			utils.DieOnError(err, "Failed to init task repo")
		}
		listRepo := &listDomain.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "list"),
		}
		projectRepo, err := repo.New(repo.Options{
			ResourceType: "projects",
			DB:           db,
			Logger:       logr.Fork("repo", "project"),
		})
		utils.DieOnError(err, "Failed to init project repo")

		userRepo := &userDomain.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "user"),
		}
		if err := userRepo.Initiate(context.Background()); err != nil {
			utils.DieOnError(err, "Failed to init user repo")
		}
		automationRepo := &automationDomain.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "automation"),
		}
		triggerRepo := &triggerDomain.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "trigger"),
		}

		ebus := internal.NewEventBusFromFlagsOrDie(db, userRepo, true, logr.Fork("module", "eventbus"))
		defer ebus.Stop()
		logr.Info("Eventbus connected")
		cb := internal.NewCommandBusFromFlagsOrDie(userRepo, logr.Fork("module", "commandbus"))
		err = cb.Start()
		logr.Info("Commandbus connected")
		registerCommandHandlers(cb, ebus, userRepo, &taskRepo)

		registerListEventHandlers(ebus, listRepo)
		registerProjectEventHandlers(ebus, projectRepo)
		registerTriggerEventHandlers(ebus, triggerRepo)
		registerAutomationEventHandlers(ebus, automationRepo)
		registerViewEventHandlers(ebus, db, &taskRepo, listRepo, projectRepo, logr)
		sagaEventHandler := &saga.EventHandler{
			CommandBus:     cb,
			TaskRepo:       &taskRepo,
			ListRepo:       listRepo,
			AutomationRepo: automationRepo,
			ProjectRepo:    projectRepo,
			TriggerRepo:    triggerRepo,
			UserRepo:       userRepo,
		}
		registerSagas(ebus, sagaEventHandler)
		err = ebus.Start()
		utils.DieOnError(err, "Failed to start eventbus")
		s.SetReady()
		return s.Start()
	},
}

func init() {
	startCmd.AddCommand(eventHandlerServiceCmd)

	viper.BindEnv("port", "PORT")
	viper.BindEnv("postgres-url", "POSTGRES_URL")
	viper.BindEnv("rabbitmq-host", "RABBITMQ_HOST")
	viper.BindEnv("rabbitmq-port", "RABBITMQ_PORT")
	viper.BindEnv("rabbitmq-api-port", "RABBITMQ_API_PORT")
	viper.BindEnv("rabbitmq-username", "RABBITMQ_USERNAME")
	viper.BindEnv("rabbitmq-password", "RABBITMQ_PASSWORD")

	viper.SetDefault("port", "8080")

	eventHandlerServiceCmd.Flags().StringVar(&eventHandlerServiceCmdFlags.port, "port", viper.GetString("port"), "Set server port")
	eventHandlerServiceCmd.Flags().StringVar(&eventHandlerServiceCmdFlags.postgresURL, "postgres-url", viper.GetString("postgres-url"), "Connection string to postgres [$POSTGRES_URL]")
	eventHandlerServiceCmd.Flags().StringVar(&eventHandlerServiceCmdFlags.rabbitmqHost, "rabbitmq-host", viper.GetString("rabbitmq-host"), "RabbitMQ host [$RABBITMQ_HOST]")
	eventHandlerServiceCmd.Flags().StringVar(&eventHandlerServiceCmdFlags.rabbitmqPort, "rabbitmq-port", viper.GetString("rabbitmq-port"), "RabbitMQ port [$RABBITMQ_PORT]")
	eventHandlerServiceCmd.Flags().StringVar(&eventHandlerServiceCmdFlags.rabbitmqAPIPort, "rabbitmq-api-port", viper.GetString("rabbitmq-api-port"), "RabbitMQ API port [$RABBITMQ_API_PORT]")
	eventHandlerServiceCmd.Flags().StringVar(&eventHandlerServiceCmdFlags.rabbitmqUsername, "rabbitmq-username", viper.GetString("rabbitmq-username"), "RabbitMQ username [$RABBITMQ_API_PORT]")

	// Set flag value when viper has value from env var
	eventHandlerServiceCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			eventHandlerServiceCmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
	eventHandlerServiceCmd.MarkFlagRequired("postgres-url")
}
