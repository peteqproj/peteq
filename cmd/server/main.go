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
	"github.com/peteqproj/peteq/pkg/api"
	"github.com/peteqproj/peteq/pkg/api/list"
	"github.com/peteqproj/peteq/pkg/api/project"
	"github.com/peteqproj/peteq/pkg/api/task"
	"github.com/peteqproj/peteq/pkg/api/view"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/db/local"
	eventbus "github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/server"
)

var locatDBLocation = path.Join(os.Getenv("HOME"), ".peteq")

func main() {
	s := server.New(server.Options{
		Port: "8080",
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
	taskRepo := &taskDomain.Repo{
		DB: taskLocalDB,
	}

	listRepo := &listDomain.Repo{
		DB: listLocalDB,
	}

	projectRepo := &projectDomain.Repo{
		DB: projectLocalDB,
	}

	cb := commandbus.New(commandbus.Options{
		Type: "local",
	})

	registerTaskEventHandlers(inmemoryEventbus, taskRepo)
	registerListEventHandlers(inmemoryEventbus, listRepo)
	registerProjectEventHandlers(inmemoryEventbus, projectRepo)
	registerCommandHandlers(cb, inmemoryEventbus)

	commandAPI := task.CommandAPI{
		Repo: &taskDomain.Repo{
			DB: taskLocalDB,
		},
		Commandbus: cb,
	}
	queryAPI := task.QueryAPI{
		Repo: &taskDomain.Repo{
			DB: taskLocalDB,
		},
	}

	projectCommandAPI := project.CommandAPI{
		Repo:       projectRepo,
		Commandbus: cb,
	}
	projectQueryAPI := project.QueryAPI{
		Repo: projectRepo,
	}

	listCommandAPI := list.CommandAPI{
		Repo:       listRepo,
		Commandbus: cb,
	}
	listQueryAPI := list.QueryAPI{
		Repo: listRepo,
	}

	backlogViewAPI := view.BacklogViewAPI{
		TaskRepo:    taskRepo,
		ListRepo:    listRepo,
		ProjectRepo: projectRepo,
	}

	projectsViewAPI := view.ProjectsViewAPI{
		TaskRepo:    taskRepo,
		ProjectRepo: projectRepo,
	}

	s.AddResource(buildAPI(buildAPIOptions{
		taskQueryAPI:      queryAPI,
		taskCommandAPI:    commandAPI,
		listQuestAPI:      listQueryAPI,
		listCommandAPI:    listCommandAPI,
		projectCommandAPI: projectCommandAPI,
		projectQueryAPI:   projectQueryAPI,
		backlogViewAPI:    backlogViewAPI,
		projectsViewAPI:   projectsViewAPI,
	}))

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
}

type buildAPIOptions struct {
	taskQueryAPI      task.QueryAPI
	taskCommandAPI    task.CommandAPI
	listQuestAPI      list.QueryAPI
	listCommandAPI    list.CommandAPI
	projectQueryAPI   project.QueryAPI
	projectCommandAPI project.CommandAPI
	backlogViewAPI    view.BacklogViewAPI
	projectsViewAPI   view.ProjectsViewAPI
}

func buildAPI(options buildAPIOptions) api.Resource {
	return api.Resource{
		Path: "/api",
		Subresource: []api.Resource{
			{
				Path: "/task",
				Endpoints: []api.Endpoint{
					{
						Verb:    "GET",
						Handler: options.taskQueryAPI.List,
					},
					{
						Verb:    "POST",
						Path:    "/complete",
						Handler: api.WrapCommandAPI(options.taskCommandAPI.Complete),
					},
					{
						Verb:    "POST",
						Path:    "/reopen",
						Handler: api.WrapCommandAPI(options.taskCommandAPI.Reopen),
					},
					{
						Verb:    "POST",
						Path:    "/create",
						Handler: api.WrapCommandAPI(options.taskCommandAPI.Create),
					},
					{
						Verb:    "POST",
						Path:    "/update",
						Handler: api.WrapCommandAPI(options.taskCommandAPI.Update),
					},
					{
						Verb:    "POST",
						Path:    "/delete",
						Handler: api.WrapCommandAPI(options.taskCommandAPI.Delete),
					},
				},
				Subresource: []api.Resource{
					{
						Path: "/:id",
						Endpoints: []api.Endpoint{
							{
								Verb:    "GET",
								Handler: options.taskQueryAPI.Get,
							},
						},
					},
				},
			},
			{
				Path: "/list",
				Endpoints: []api.Endpoint{
					{
						Verb:    "GET",
						Handler: options.listQuestAPI.List,
					},
					{
						Verb:    "POST",
						Path:    "/moveTasks",
						Handler: api.WrapCommandAPI(options.listCommandAPI.MoveTasks),
					},
				},
			},
			{
				Path: "/project",
				Endpoints: []api.Endpoint{
					{
						Verb:    "GET",
						Handler: options.projectQueryAPI.List,
					},
					{
						Path:    "/create",
						Verb:    "POST",
						Handler: api.WrapCommandAPI(options.projectCommandAPI.Create),
					},
					{
						Path:    "/addTasks",
						Verb:    "POST",
						Handler: api.WrapCommandAPI(options.projectCommandAPI.AddTasks),
					},
				},
				Subresource: []api.Resource{
					{
						Path: "/:id",
						Endpoints: []api.Endpoint{
							{
								Verb:    "GET",
								Handler: options.projectQueryAPI.Get,
							},
						},
					},
				},
			},
			{
				Path: "/view",
				Endpoints: []api.Endpoint{
					{
						Path:    "/backlog",
						Handler: options.backlogViewAPI.Get,
						Verb:    "GET",
					},
					{
						Path:    "/projects",
						Handler: options.projectsViewAPI.Get,
						Verb:    "GET",
					},
				},
			},
		},
	}
}
