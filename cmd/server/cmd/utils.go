package cmd

import (
	_ "github.com/lib/pq"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/peteqproj/peteq/saga"

	"github.com/peteqproj/peteq/domain/automation"
	automationCommands "github.com/peteqproj/peteq/domain/automation/command"
	"github.com/peteqproj/peteq/domain/list"
	listDomain "github.com/peteqproj/peteq/domain/list"
	listCommands "github.com/peteqproj/peteq/domain/list/command"
	"github.com/peteqproj/peteq/domain/project"
	projectDomain "github.com/peteqproj/peteq/domain/project"
	projectCommands "github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/domain/sensor"
	sensorCommands "github.com/peteqproj/peteq/domain/sensor/command"
	sensorEventTypes "github.com/peteqproj/peteq/domain/sensor/event/types"
	"github.com/peteqproj/peteq/domain/task"
	taskCommands "github.com/peteqproj/peteq/domain/task/command"
	userDomain "github.com/peteqproj/peteq/domain/user"
	userCommands "github.com/peteqproj/peteq/domain/user/command"
	userEventTypes "github.com/peteqproj/peteq/domain/user/event/types"
	viewBuilder "github.com/peteqproj/peteq/pkg/api/view/builder"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	eventbus "github.com/peteqproj/peteq/pkg/event/bus"
)

// DieOnError kills the process and prints a message
func DieOnError(err error, msg string) {
	utils.DieOnError(err, msg)
}

func registerCommandHandlers(cb commandbus.CommandBus, eventbus eventbus.EventPublisher, userRepo *userDomain.Repo, taskRepo *task.Repo, listRepo *list.Repo, projectRepo *project.Repo, sensorRepo *sensor.Repo, automationRepo *automation.Repo) {
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
		Repo:     listRepo,
	})
	cb.RegisterHandler("list.create", &listCommands.Create{
		Eventbus: eventbus,
		Repo:     listRepo,
	})

	// Project related commands
	cb.RegisterHandler("project.create", &projectCommands.CreateCommand{
		Eventbus: eventbus,
		Repo:     projectRepo,
	})
	cb.RegisterHandler("project.add-task", &projectCommands.AddTaskCommand{
		Eventbus: eventbus,
		Repo:     projectRepo,
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

	// Sensor related commands
	cb.RegisterHandler("sensor.create", &sensorCommands.CreateCommand{
		Eventbus: eventbus,
		Repo:     sensorRepo,
	})
	cb.RegisterHandler("sensor.trigger", &sensorCommands.TriggerCommand{
		Eventbus: eventbus,
		Repo:     sensorRepo,
	})

	// Automation related commands
	cb.RegisterHandler("automation.create", &automationCommands.CreateCommand{
		Eventbus: eventbus,
		Repo:     automationRepo,
	})
	cb.RegisterHandler("automation.bindSensor", &automationCommands.CreateSensorBindingCommand{
		Eventbus: eventbus,
		Repo:     automationRepo,
	})
}

func registerSagas(eventbus eventbus.Eventbus, eh *saga.EventHandler) {
	eventbus.Subscribe(sensorEventTypes.SensorTriggeredEvent, eh)
	eventbus.Subscribe(userEventTypes.UserRegistredEvent, eh)
}

func registerViewEventHandlers(eventbus eventbus.Eventbus, db db.Database, taskRepo *task.Repo, listRepo *listDomain.Repo, projectRepo *projectDomain.Repo, logger logger.Logger) {
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
