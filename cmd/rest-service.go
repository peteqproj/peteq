package cmd

import (
	listDomain "github.com/peteqproj/peteq/domain/list"
	projectDomain "github.com/peteqproj/peteq/domain/project"
	taskDomain "github.com/peteqproj/peteq/domain/task"
	userDomain "github.com/peteqproj/peteq/domain/user"
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

var restServiceCmdFlags struct {
	verbose          bool
	postgresURL      string
	rabbitmqHost     string
	rabbitmqPort     string
	rabbitmqAPIPort  string
	rabbitmqUsername string
	rabbitmqPassword string
	port             string
}

var restServiceCmd = &cobra.Command{
	Use:   "rest-api-service",
	Short: "Starts the RESTful API service",
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

		apiBuilder := builder.Builder{
			UserRepo:    userRepo,
			ListRpeo:    listRepo,
			ProjectRepo: projectRepo,
			TaskRepo:    taskRepo,
			DB:          db,
			Logger:      logr,
		}
		s.AddResource(apiBuilder.BuildRestfulAPI())
		s.SetReady()
		return s.Start()
	},
}

func init() {
	startCmd.AddCommand(restServiceCmd)

	viper.BindEnv("port", "PORT")
	viper.BindEnv("postgres-url", "POSTGRES_URL")
	viper.BindEnv("rabbitmq-host", "RABBITMQ_HOST")
	viper.BindEnv("rabbitmq-port", "RABBITMQ_PORT")
	viper.BindEnv("rabbitmq-api-port", "RABBITMQ_API_PORT")
	viper.BindEnv("rabbitmq-username", "RABBITMQ_USERNAME")
	viper.BindEnv("rabbitmq-password", "RABBITMQ_PASSWORD")

	viper.SetDefault("port", "8080")

	restServiceCmd.Flags().StringVar(&restServiceCmdFlags.port, "port", viper.GetString("port"), "Set server port")
	restServiceCmd.Flags().StringVar(&restServiceCmdFlags.postgresURL, "postgres-url", viper.GetString("postgres-url"), "Connection string to postgres [$POSTGRES_URL]")
	restServiceCmd.Flags().StringVar(&restServiceCmdFlags.rabbitmqHost, "rabbitmq-host", viper.GetString("rabbitmq-host"), "RabbitMQ host [$RABBITMQ_HOST]")
	restServiceCmd.Flags().StringVar(&restServiceCmdFlags.rabbitmqPort, "rabbitmq-port", viper.GetString("rabbitmq-port"), "RabbitMQ port [$RABBITMQ_PORT]")
	restServiceCmd.Flags().StringVar(&restServiceCmdFlags.rabbitmqAPIPort, "rabbitmq-api-port", viper.GetString("rabbitmq-api-port"), "RabbitMQ API port [$RABBITMQ_API_PORT]")
	restServiceCmd.Flags().StringVar(&restServiceCmdFlags.rabbitmqUsername, "rabbitmq-username", viper.GetString("rabbitmq-username"), "RabbitMQ username [$RABBITMQ_API_PORT]")

	// Set flag value when viper has value from env var
	restServiceCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			restServiceCmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
	restServiceCmd.MarkFlagRequired("postgres-url")
}
