package main

import (
	"fmt"

	_ "github.com/lib/pq"

	socketio "github.com/googollee/go-socket.io"

	listDomain "github.com/peteqproj/peteq/domain/list"
	listCommands "github.com/peteqproj/peteq/domain/list/command"
	listEventHandlers "github.com/peteqproj/peteq/domain/list/event/handler"
	projectDomain "github.com/peteqproj/peteq/domain/project"
	projectCommands "github.com/peteqproj/peteq/domain/project/command"
	projectEventHandlers "github.com/peteqproj/peteq/domain/project/event/handler"
	taskDomain "github.com/peteqproj/peteq/domain/task"
	taskCommands "github.com/peteqproj/peteq/domain/task/command"
	taskEventHandlers "github.com/peteqproj/peteq/domain/task/event/handler"
	userDomain "github.com/peteqproj/peteq/domain/user"
	userCommands "github.com/peteqproj/peteq/domain/user/command"
	userEventHandlers "github.com/peteqproj/peteq/domain/user/event/handler"
	"github.com/peteqproj/peteq/pkg/api/builder"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/config"
	"github.com/peteqproj/peteq/pkg/db/postgres"
	eventbus "github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/server"
	"github.com/peteqproj/peteq/pkg/utils"
)

func main() {
	logr := logger.New(logger.Options{})
	cnf := &config.Server{
		Port:                 utils.GetEnvOrDie("PORT"),
		EncryptionPassphrase: "local",
	}
	s := server.New(server.Options{
		Config: cnf,
	})
	wsserver, err := socketio.NewServer(nil)
	defer wsserver.Close()
	utils.DieOnError(err, "Failed to create WS server")
	handlerWSEvents(wsserver)
	err = s.AddWS(wsserver)
	utils.DieOnError(err, "Failed to attach WS server")

	db, err := postgres.Connect(utils.GetEnvOrDie("POSTGRES_URL"))
	defer db.Close()
	utils.DieOnError(err, "Failed to connect to postgres")

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

	cb := commandbus.New(commandbus.Options{
		Type:   "local",
		Logger: logr.Fork("module", "commandbus"),
	})

	registerUserEventHandlers(ebus, userRepo)
	registerTaskEventHandlers(ebus, taskRepo)
	registerListEventHandlers(ebus, listRepo)
	registerProjectEventHandlers(ebus, projectRepo)
	registerCommandHandlers(cb, ebus)

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

	s.AddResource(apiBuilder.BuildCommandAPI())
	s.AddResource(apiBuilder.BuildViewAPI())
	s.AddResource(apiBuilder.BuildRestfulAPI())
	err = ebus.Start()
	utils.DieOnError(err, "Failed to start eventbus")
	err = s.Start()
	utils.DieOnError(err, "Failed to run server")
}

func handlerWSEvents(server *socketio.Server) {
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println(msg)
		s.Emit("reply", "have "+msg)
	})
	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})
	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})

	go server.Serve()
}

func registerTaskEventHandlers(eventbus eventbus.Eventbus, repo *taskDomain.Repo) {
	// Task related event handlers
	eventbus.Subscribe("task.created", &taskEventHandlers.CreatedHandler{
		Repo: repo,
	})
	eventbus.Subscribe("task.deleted", &taskEventHandlers.DeleteHandler{
		Repo: repo,
	})
	eventbus.Subscribe("task.updated", &taskEventHandlers.UpdatedHandler{
		Repo: repo,
	})
	eventbus.Subscribe("task.completed", &taskEventHandlers.CompletedHandler{
		Repo: repo,
	})
	eventbus.Subscribe("task.reopened", &taskEventHandlers.ReopenedHandler{
		Repo: repo,
	})
}

func registerListEventHandlers(eventbus eventbus.Eventbus, repo *listDomain.Repo) {
	// List related event handlers
	eventbus.Subscribe("list.task-moved", &listEventHandlers.TaskMovedHandler{
		Repo: repo,
	})
	eventbus.Subscribe("list.created", &listEventHandlers.CreatedHandler{
		Repo: repo,
	})
}

func registerUserEventHandlers(eventbus eventbus.Eventbus, repo *userDomain.Repo) {
	// User related event handlers
	eventbus.Subscribe("user.registred", &userEventHandlers.RegistredHandler{
		Repo: repo,
	})
	eventbus.Subscribe("user.loggedin", &userEventHandlers.LoggedinHandler{
		Repo: repo,
	})
}

func registerProjectEventHandlers(eventbus eventbus.Eventbus, repo *projectDomain.Repo) {
	// List related event handlers
	eventbus.Subscribe("project.created", &projectEventHandlers.CreatedHandler{
		Repo: repo,
	})

	eventbus.Subscribe("project.task-added", &projectEventHandlers.TaskAddedHandler{
		Repo: repo,
	})
}

func registerCommandHandlers(cb commandbus.CommandBus, eventbus eventbus.Eventbus) {
	// Task related commands
	cb.RegisterHandler("task.create", &taskCommands.CreateCommand{
		Eventbus: eventbus,
	})
	cb.RegisterHandler("task.delete", &taskCommands.DeleteCommand{
		Eventbus: eventbus,
	})
	cb.RegisterHandler("task.update", &taskCommands.UpdateCommand{
		Eventbus: eventbus,
	})
	cb.RegisterHandler("task.complete", &taskCommands.CompleteCommand{
		Eventbus: eventbus,
	})
	cb.RegisterHandler("task.reopen", &taskCommands.ReopenCommand{
		Eventbus: eventbus,
	})

	// List related command
	cb.RegisterHandler("list.move-task", &listCommands.MoveTaskCommand{
		Eventbus: eventbus,
	})
	cb.RegisterHandler("list.create", &listCommands.Create{
		Eventbus: eventbus,
	})

	// Project related commands
	cb.RegisterHandler("project.create", &projectCommands.CreateCommand{
		Eventbus: eventbus,
	})
	cb.RegisterHandler("project.add-task", &projectCommands.AddTaskCommand{
		Eventbus: eventbus,
	})

	// User related commands
	cb.RegisterHandler("user.register", &userCommands.RegisterCommand{
		Eventbus: eventbus,
	})
	cb.RegisterHandler("user.login", &userCommands.LoginCommand{
		Eventbus: eventbus,
	})
}
