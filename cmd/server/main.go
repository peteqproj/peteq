package main

import (
	_ "github.com/lib/pq"

	automationDomain "github.com/peteqproj/peteq/domain/automation"
	automationCommands "github.com/peteqproj/peteq/domain/automation/command"
	automationEventHandlers "github.com/peteqproj/peteq/domain/automation/event/handler"
	automationEventTypes "github.com/peteqproj/peteq/domain/automation/event/types"
	listDomain "github.com/peteqproj/peteq/domain/list"
	listCommands "github.com/peteqproj/peteq/domain/list/command"
	listEventHandlers "github.com/peteqproj/peteq/domain/list/event/handler"
	listEventTypes "github.com/peteqproj/peteq/domain/list/event/types"
	projectDomain "github.com/peteqproj/peteq/domain/project"
	projectCommands "github.com/peteqproj/peteq/domain/project/command"
	projectEventHandlers "github.com/peteqproj/peteq/domain/project/event/handler"
	projectEventTypes "github.com/peteqproj/peteq/domain/project/event/types"
	taskDomain "github.com/peteqproj/peteq/domain/task"
	taskCommands "github.com/peteqproj/peteq/domain/task/command"
	taskEventHandlers "github.com/peteqproj/peteq/domain/task/event/handler"
	triggerDomain "github.com/peteqproj/peteq/domain/trigger"
	triggerCommands "github.com/peteqproj/peteq/domain/trigger/command"
	triggerEventHandlers "github.com/peteqproj/peteq/domain/trigger/event/handler"
	triggerEventTypes "github.com/peteqproj/peteq/domain/trigger/event/types"
	userDomain "github.com/peteqproj/peteq/domain/user"
	userCommands "github.com/peteqproj/peteq/domain/user/command"
	userEventHandlers "github.com/peteqproj/peteq/domain/user/event/handler"
	userEventTypes "github.com/peteqproj/peteq/domain/user/event/types"
	"github.com/peteqproj/peteq/internal"
	"github.com/peteqproj/peteq/pkg/api/builder"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/config"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/db/postgres"
	eventbus "github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/server"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/peteqproj/peteq/saga"

	taskEventTypes "github.com/peteqproj/peteq/domain/task/event/types"
)

type (
	sagaEventHandler struct {
		cb       commandbus.CommandBus
		taskRepo *taskDomain.Repo
		listRepo *listDomain.Repo
		lgr      logger.Logger
	}
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
	automationRepo := &automationDomain.Repo{
		DB:     db,
		Logger: logr.Fork("repo", "automation"),
	}

	userRepo := &userDomain.Repo{
		DB:     db,
		Logger: logr.Fork("repo", "user"),
	}

	triggerRepo := &triggerDomain.Repo{
		DB:     db,
		Logger: logr.Fork("repo", "trigger"),
	}

	ebus := internal.NewEventBusFromFlagsOrDie(db, userRepo, logr.Fork("module", "eventbus"))
	defer ebus.Stop()
	logr.Info("Eventbus connected")
	cb := internal.NewCommandBusFromFlagsOrDie(userRepo, logr.Fork("module", "commandbus"))
	err = cb.Start()
	logr.Info("Commandbus connected")

	registerUserEventHandlers(ebus, userRepo)
	registerTaskEventHandlers(ebus, taskRepo)
	registerListEventHandlers(ebus, listRepo)
	registerProjectEventHandlers(ebus, projectRepo)
	registerTriggerEventHandlers(ebus, triggerRepo)
	registerAutomationEventHandlers(ebus, automationRepo)
	registerCommandHandlers(cb, ebus, userRepo)

	sagaEventHandler := &saga.EventHandler{
		CommandBus:     cb,
		TaskRepo:       taskRepo,
		ListRepo:       listRepo,
		AutomationRepo: automationRepo,
		ProjectRepo:    projectRepo,
		TriggerRepo:    triggerRepo,
		UserRepo:       userRepo,
	}
	registerSagas(ebus, sagaEventHandler)

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

func registerTaskEventHandlers(eventbus eventbus.Eventbus, repo *taskDomain.Repo) {
	// Task related event handlers
	eventbus.Subscribe(taskEventTypes.TaskCreatedEvent, &taskEventHandlers.CreatedHandler{
		Repo: repo,
	})
	eventbus.Subscribe(taskEventTypes.TaskDeletedEvent, &taskEventHandlers.DeleteHandler{
		Repo: repo,
	})
	eventbus.Subscribe(taskEventTypes.TaskUpdatedEvent, &taskEventHandlers.UpdatedHandler{
		Repo: repo,
	})
	eventbus.Subscribe(taskEventTypes.TaskStatusChanged, &taskEventHandlers.StatusChangedHandler{
		Repo: repo,
	})
	eventbus.Subscribe(taskEventTypes.TaskStatusChanged, &taskEventHandlers.StatusChangedHandler{
		Repo: repo,
	})
}

func registerListEventHandlers(eventbus eventbus.Eventbus, repo *listDomain.Repo) {
	// List related event handlers
	eventbus.Subscribe(listEventTypes.TaskMovedIntoListEvent, &listEventHandlers.TaskMovedHandler{
		Repo: repo,
	})
	eventbus.Subscribe(listEventTypes.ListCreatedEvent, &listEventHandlers.CreatedHandler{
		Repo: repo,
	})
}

func registerUserEventHandlers(eventbus eventbus.Eventbus, repo *userDomain.Repo) {
	// User related event handlers
	eventbus.Subscribe(userEventTypes.UserRegistredEvent, &userEventHandlers.RegistredHandler{
		Repo: repo,
	})
	eventbus.Subscribe(userEventTypes.UserLoggedIn, &userEventHandlers.LoggedinHandler{
		Repo: repo,
	})
}

func registerProjectEventHandlers(eventbus eventbus.Eventbus, repo *projectDomain.Repo) {
	// List related event handlers
	eventbus.Subscribe(projectEventTypes.ProjectCreatedEvent, &projectEventHandlers.CreatedHandler{
		Repo: repo,
	})

	eventbus.Subscribe(projectEventTypes.TaskAddedToProjectEvent, &projectEventHandlers.TaskAddedHandler{
		Repo: repo,
	})
}

func registerTriggerEventHandlers(eventbus eventbus.Eventbus, repo *triggerDomain.Repo) {
	// Trigger related event handlers
	eventbus.Subscribe(triggerEventTypes.TriggerCreatedEvent, &triggerEventHandlers.CreatedHandler{
		Repo: repo,
	})
}

func registerAutomationEventHandlers(eventbus eventbus.Eventbus, repo *automationDomain.Repo) {
	// Automation related event handlers
	eventbus.Subscribe(automationEventTypes.AutomationCreatedEvent, &automationEventHandlers.CreatedHandler{
		Repo: repo,
	})
	eventbus.Subscribe(automationEventTypes.TriggerBindingCreatedEvent, &automationEventHandlers.TriggerBindingCreatedHandler{
		Repo: repo,
	})
}

func registerCommandHandlers(cb commandbus.CommandBus, eventbus eventbus.Eventbus, userRepo *userDomain.Repo) {
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
		Eventbus:    eventbus,
		Commandbus:  cb,
		IDGenerator: utils.NewGenerator(),
		Repo:        userRepo,
	})
	cb.RegisterHandler("user.login", &userCommands.LoginCommand{
		Eventbus: eventbus,
		Repo:     userRepo,
	})

	// Trigger related commands
	cb.RegisterHandler("trigger.create", &triggerCommands.CreateCommand{
		Eventbus: eventbus,
	})

	// Automation related commands
	cb.RegisterHandler("automation.create", &automationCommands.CreateCommand{
		Eventbus: eventbus,
	})
	cb.RegisterHandler("automation.bindTrigger", &automationCommands.CreateTriggerBindingCommand{
		Eventbus: eventbus,
	})
}

func registerSagas(eventbus eventbus.Eventbus, eh *saga.EventHandler) {
	eventbus.Subscribe(triggerEventTypes.TriggerTriggeredEvent, eh)
	eventbus.Subscribe(userEventTypes.UserRegistredEvent, eh)
}
