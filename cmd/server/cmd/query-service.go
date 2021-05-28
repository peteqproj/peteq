package cmd

import (
	"context"

	listDomain "github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/task"
	userDomain "github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/pkg/api/builder"
	"github.com/peteqproj/peteq/pkg/config"
	"github.com/peteqproj/peteq/pkg/db"
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

		db, err := db.New(db.Options{
			URL: utils.GetEnvOrDie("POSTGRES_URL"),
		})
		utils.DieOnError(err, "Failed to connect to postgres")
		sqldb, err := db.DB()
		utils.DieOnError(err, "Failed to get sql object")
		defer sqldb.Close()
		taskRepo := &task.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "task"),
		}
		if err := taskRepo.Initiate(context.Background()); err != nil {
			utils.DieOnError(err, "Failed to init task repo")
		}
		utils.DieOnError(err, "Failed to init task repo")
		listRepo := &listDomain.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "list"),
		}
		if err := taskRepo.Initiate(context.Background()); err != nil {
			utils.DieOnError(err, "Failed to init list repo")
		}
		projectRepo := &project.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "project"),
		}
		if err := projectRepo.Initiate(context.Background()); err != nil {
			utils.DieOnError(err, "Failed to init project repo")
		}

		userRepo := &userDomain.Repo{
			DB:     db,
			Logger: logr.Fork("repo", "user"),
		}
		if err := userRepo.Initiate(context.Background()); err != nil {
			utils.DieOnError(err, "Failed to init user repo")
		}
		apiBuilder := builder.Builder{
			UserRepo:    userRepo,
			ListRpeo:    listRepo,
			ProjectRepo: projectRepo,
			TaskRepo:    taskRepo,
			DB:          db,
			Logger:      logr,
		}
		s.AddResource(apiBuilder.BuildViewAPI())
		s.SetReady()
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
