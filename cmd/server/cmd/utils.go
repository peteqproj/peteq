package cmd

import (
	_ "github.com/lib/pq"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/repo"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/peteqproj/peteq/saga"

	automationDomain "github.com/peteqproj/peteq/domain/automation"
	automationCommands "github.com/peteqproj/peteq/domain/automation/command"
	automationEventHandlers "github.com/peteqproj/peteq/domain/automation/event/handler"
	automationEventTypes "github.com/peteqproj/peteq/domain/automation/event/types"
	listDomain "github.com/peteqproj/peteq/domain/list"
	listCommands "github.com/peteqproj/peteq/domain/list/command"
	listEventHandlers "github.com/peteqproj/peteq/domain/list/event/handler"
	listEventTypes "github.com/peteqproj/peteq/domain/list/event/types"
	projectCommands "github.com/peteqproj/peteq/domain/project/command"
	projectEventHandlers "github.com/peteqproj/peteq/domain/project/event/handler"
	projectEventTypes "github.com/peteqproj/peteq/domain/project/event/types"
	"github.com/peteqproj/peteq/domain/task"
	taskCommands "github.com/peteqproj/peteq/domain/task/command"
	triggerDomain "github.com/peteqproj/peteq/domain/trigger"
	triggerCommands "github.com/peteqproj/peteq/domain/trigger/command"
	triggerEventHandlers "github.com/peteqproj/peteq/domain/trigger/event/handler"
	triggerEventTypes "github.com/peteqproj/peteq/domain/trigger/event/types"
	userDomain "github.com/peteqproj/peteq/domain/user"
	userCommands "github.com/peteqproj/peteq/domain/user/command"
	userEventHandlers "github.com/peteqproj/peteq/domain/user/event/handler"
	userEventTypes "github.com/peteqproj/peteq/domain/user/event/types"
	viewBuilder "github.com/peteqproj/peteq/pkg/api/view/builder"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	eventbus "github.com/peteqproj/peteq/pkg/event/bus"
)

// DieOnError kills the process and prints a message
func DieOnError(err error, msg string) {
	utils.DieOnError(err, msg)
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

func registerProjectEventHandlers(eventbus eventbus.Eventbus, repo *repo.Repo) {
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

func registerCommandHandlers(cb commandbus.CommandBus, eventbus eventbus.EventPublisher, userRepo *userDomain.Repo, taskRepo *task.Repo) {
	// Task related commands
	cb.RegisterHandler("task.create", &taskCommands.CreateCommand{
		Eventbus: eventbus,
		Repo:     taskRepo,
	})
	cb.RegisterHandler("task.delete", &taskCommands.DeleteCommand{
		Eventbus: eventbus,
		Repo:     taskRepo,
	})
	cb.RegisterHandler("task.update", &taskCommands.UpdateCommand{
		Eventbus: eventbus,
		Repo:     taskRepo,
	})
	cb.RegisterHandler("task.complete", &taskCommands.CompleteCommand{
		Eventbus: eventbus,
		Repo:     taskRepo,
	})
	cb.RegisterHandler("task.reopen", &taskCommands.ReopenCommand{
		Eventbus: eventbus,
		Repo:     taskRepo,
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
	cb.RegisterHandler("trigger.run", &triggerCommands.RunCommand{
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

func registerViewEventHandlers(eventbus eventbus.Eventbus, db db.Database, taskRepo *task.Repo, listRepo *listDomain.Repo, projectRepo *repo.Repo, logger logger.Logger) {
	vb := viewBuilder.New(&viewBuilder.Options{
		TaskRepo:    taskRepo,
		ListRepo:    listRepo,
		ProjectRepo: projectRepo,
		Logger:      logger,
		DB:          db,
	})
	views := vb.BuildViews()
	for _, view := range views {
		for name, handler := range view.EventHandlers() {
			logger.Info("Subscribing", "name", name, "handler", handler.Name())
			eventbus.Subscribe(name, handler)
		}
	}
}
