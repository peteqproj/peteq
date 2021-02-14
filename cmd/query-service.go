package cmd

import (
	_ "github.com/lib/pq"

	listDomain "github.com/peteqproj/peteq/domain/list"
	projectDomain "github.com/peteqproj/peteq/domain/project"
	taskDomain "github.com/peteqproj/peteq/domain/task"
	userDomain "github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/internal"
	"github.com/peteqproj/peteq/pkg/api/builder"
	"github.com/peteqproj/peteq/pkg/config"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/db/postgres"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/server"
	"github.com/peteqproj/peteq/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var queryServiceCmdFlags struct {
	verbose          bool
	postgresURL      string
	rabbitmqHost     string
	rabbitmqPort     string
	rabbitmqAPIPort  string
	rabbitmqUsername string
	rabbitmqPassword string
	port             string
}

var queryServiceCmd = &cobra.Command{
	Use:   "query-service",
	Short: "Starts the query API service (The Q of CQRS)",
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

		taskRepo := &taskDomain.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "task"),
		}
		listRepo := &listDomain.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "list"),
		}
		projectRepo := &projectDomain.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "project"),
		}
		userRepo := &userDomain.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "user"),
		}

		ebus := internal.NewEventBusFromFlagsOrDie(db, userRepo, logr.Fork("module", "eventbus"))
		defer ebus.Stop()
		logr.Info("Eventbus connected")
		cb := internal.NewCommandBusFromFlagsOrDie(userRepo, logr.Fork("module", "commandbus"))
		err = cb.Start()
		logr.Info("Commandbus connected")
		registerCommandHandlers(cb, ebus, userRepo)

		apiBuilder := builder.Builder{
			UserRepo:    userRepo,
			ListRpeo:    listRepo,
			ProjectRepo: projectRepo,
			TaskRepo:    taskRepo,
			Commandbus:  cb,
			DB:          db,
			Eventbus:    ebus,
			Logger:      logr,
		}
		s.AddResource(apiBuilder.BuildViewAPI())
		err = ebus.Start()
		utils.DieOnError(err, "Failed to start eventbus")
		return s.Start()
	},
}

func init() {
	startCmd.AddCommand(queryServiceCmd)

	viper.BindEnv("port", "PORT")
	viper.BindEnv("postgres-url", "POSTGRES_URL")
	viper.BindEnv("rabbitmq-host", "RABBITMQ_HOST")
	viper.BindEnv("rabbitmq-port", "RABBITMQ_PORT")
	viper.BindEnv("rabbitmq-api-port", "RABBITMQ_API_PORT")
	viper.BindEnv("rabbitmq-username", "RABBITMQ_USERNAME")
	viper.BindEnv("rabbitmq-password", "RABBITMQ_PASSWORD")

	viper.SetDefault("port", "8080")

	queryServiceCmd.Flags().StringVar(&queryServiceCmdFlags.port, "port", viper.GetString("port"), "Set server port")
	queryServiceCmd.Flags().StringVar(&queryServiceCmdFlags.postgresURL, "postgres-url", viper.GetString("postgres-url"), "Connection string to postgres [$POSTGRES_URL]")
	queryServiceCmd.Flags().StringVar(&queryServiceCmdFlags.rabbitmqHost, "rabbitmq-host", viper.GetString("rabbitmq-host"), "RabbitMQ host [$RABBITMQ_HOST]")
	queryServiceCmd.Flags().StringVar(&queryServiceCmdFlags.rabbitmqPort, "rabbitmq-port", viper.GetString("rabbitmq-port"), "RabbitMQ port [$RABBITMQ_PORT]")
	queryServiceCmd.Flags().StringVar(&queryServiceCmdFlags.rabbitmqAPIPort, "rabbitmq-api-port", viper.GetString("rabbitmq-api-port"), "RabbitMQ API port [$RABBITMQ_API_PORT]")
	queryServiceCmd.Flags().StringVar(&queryServiceCmdFlags.rabbitmqUsername, "rabbitmq-username", viper.GetString("rabbitmq-username"), "RabbitMQ username [$RABBITMQ_API_PORT]")

	// Set flag value when viper has value from env var
	queryServiceCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			queryServiceCmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
	queryServiceCmd.MarkFlagRequired("postgres-url")
}
