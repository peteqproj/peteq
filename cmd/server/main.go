package main

import (
	"fmt"
	"os"
	"path"

	"github.com/gofrs/uuid"
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
	"github.com/peteqproj/peteq/pkg/api/apis/list"
	"github.com/peteqproj/peteq/pkg/api/builder"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/config"
	"github.com/peteqproj/peteq/pkg/db/local"
	eventbus "github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/server"
)

var locatDBLocation = path.Join(os.Getenv("HOME"), ".peteq")

func main() {
	logr := logger.New(logger.Options{})
	cnf := &config.Server{
		Port:                 "8080",
		EncryptionPassphrase: "local",
	}
	s := server.New(server.Options{
		Config: cnf,
	})
	wsserver, err := socketio.NewServer(nil)
	defer wsserver.Close()
	dieOnError(err, "Failed to create WS server")
	handlerWSEvents(wsserver)
	err = s.AddWS(wsserver)
	dieOnError(err, "Failed to attach WS server")

	taskEventStore := &local.DB{
		Path: path.Join(locatDBLocation, "tasks-events.yaml"),
	}

	inmemoryEventbus := eventbus.New(eventbus.Options{
		Type:            "local",
		LocalEventStore: taskEventStore,
		WS:              wsserver,
		Logger:          logr.Fork("module", "eventbus"),
	})
	taskLocalDB := &local.DB{
		Path: path.Join(locatDBLocation, "tasks.yaml"),
	}

	listLocalDB := &local.DB{
		Path: path.Join(locatDBLocation, "lists.yaml"),
	}

	projectLocalDB := &local.DB{
		Path: path.Join(locatDBLocation, "projects.yaml"),
	}

	userLocalDB := &local.DB{
		Path: path.Join(locatDBLocation, "users.yaml"),
	}
	taskRepo := &taskDomain.Repo{
		DB:     taskLocalDB,
		Logger: logr.Fork("repo", "task"),
	}

	listRepo := &listDomain.Repo{
		DB:     listLocalDB,
		Logger: logr.Fork("repo", "list"),
	}

	projectRepo := &projectDomain.Repo{
		DB:     projectLocalDB,
		Logger: logr.Fork("repo", "project"),
	}

	userRepo := &userDomain.Repo{
		DB:     userLocalDB,
		Logger: logr.Fork("repo", "user"),
	}

	cb := commandbus.New(commandbus.Options{
		Type:   "local",
		Logger: logr.Fork("module", "commandbus"),
	})

	registerTaskEventHandlers(inmemoryEventbus, taskRepo)
	registerListEventHandlers(inmemoryEventbus, listRepo)
	registerProjectEventHandlers(inmemoryEventbus, projectRepo)
	registerUserEventHandlers(inmemoryEventbus, userRepo)
	registerCommandHandlers(cb, inmemoryEventbus)

	listQueryAPI := list.QueryAPI{
		Repo: listRepo,
	}

	apiBuilder := builder.Builder{
		UserRepo:    userRepo,
		ListRpeo:    listRepo,
		ProjectRepo: projectRepo,
		TaskRepo:    taskRepo,
		Commandbus:  cb,
		Logger:      logr,
	}

	s.AddResource(apiBuilder.BuildCommandAPI())
	s.AddResource(apiBuilder.BuildViewAPI())

	initiateLists(*listQueryAPI.Repo)
	err = s.Start()
	dieOnError(err, "Failed to run server")
}

func initiateLists(repo listDomain.Repo) {
	if _, err := os.Stat(repo.DB.Path); err != nil {
		if os.IsNotExist(err) {
			os.Create(repo.DB.Path)
			for _, name := range []string{"This Week", "Today", "Done"} {
				repo.Create(listDomain.List{
					Metadata: listDomain.Metadata{
						ID:   uuid.Must(uuid.NewV4()).String(),
						Name: name,
					},
				})
			}
		}
	}
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
	go eventbus.Subscribe("task.created", &taskEventHandlers.CreatedHandler{
		Repo: repo,
	})
	go eventbus.Subscribe("task.deleted", &taskEventHandlers.DeleteHandler{
		Repo: repo,
	})
	go eventbus.Subscribe("task.updated", &taskEventHandlers.UpdatedHandler{
		Repo: repo,
	})
	go eventbus.Subscribe("task.completed", &taskEventHandlers.CompletedHandler{
		Repo: repo,
	})
	go eventbus.Subscribe("task.reopened", &taskEventHandlers.ReopenedHandler{
		Repo: repo,
	})
}

func registerListEventHandlers(eventbus eventbus.Eventbus, repo *listDomain.Repo) {
	// List related event handlers
	go eventbus.Subscribe("list.task-moved", &listEventHandlers.TaskMovedHandler{
		Repo: repo,
	})
}

func registerUserEventHandlers(eventbus eventbus.Eventbus, repo *userDomain.Repo) {
	// User related event handlers
	go eventbus.Subscribe("user.registred", &userEventHandlers.RegistredHandler{
		Repo: repo,
	})
	go eventbus.Subscribe("user.loggedin", &userEventHandlers.LoggedinHandler{
		Repo: repo,
	})
}

func registerProjectEventHandlers(eventbus eventbus.Eventbus, repo *projectDomain.Repo) {
	// List related event handlers
	go eventbus.Subscribe("project.created", &projectEventHandlers.CreatedHandler{
		Repo: repo,
	})

	go eventbus.Subscribe("project.task-added", &projectEventHandlers.TaskAddedHandler{
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
