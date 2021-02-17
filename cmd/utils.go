package cmd

import (
	"context"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	_ "github.com/lib/pq"
	"github.com/peteqproj/peteq/pkg/client"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/peteqproj/peteq/saga"
	"gopkg.in/yaml.v2"

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
	viewBuilder "github.com/peteqproj/peteq/pkg/api/view/builder"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	eventbus "github.com/peteqproj/peteq/pkg/event/bus"

	taskEventTypes "github.com/peteqproj/peteq/domain/task/event/types"
)

type (
	clientConfig struct {
		URL   string `yaml:"url"`
		Token string `yaml:"token"`
	}
)

// DieOnError kills the process and prints a message
func DieOnError(err error, msg string) {
	utils.DieOnError(err, msg)
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

func registerCommandHandlers(cb commandbus.CommandBus, eventbus eventbus.EventPublisher, userRepo *userDomain.Repo) {
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

func registerViewEventHandlers(eventbus eventbus.Eventbus, db db.Database, taskRepo *taskDomain.Repo, listRepo *listDomain.Repo, projectRepo *projectDomain.Repo, logger logger.Logger) {
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

func createClientConfiguration() (*client.Configuration, context.Context, error) {
	c := &clientConfig{}
	data, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".peteq/config"))
	if err != nil {
		return nil, nil, err
	}
	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, nil, err
	}
	u, err := url.Parse(c.URL)
	if err != nil {
		return nil, nil, err
	}

	cnf := &client.Configuration{
		DefaultHeader: make(map[string]string),
		UserAgent:     "peteq-cli",
		Debug:         false,
		Scheme:        u.Scheme,
		Servers: client.ServerConfigurations{
			{
				URL: u.Host,
			},
		},
	}
	ctx := context.WithValue(context.Background(), client.ContextAPIKeys, map[string]client.APIKey{
		"ApiKeyAuth": {
			Key: c.Token,
		},
	})
	return cnf, ctx, nil
}

func storeClientConfiguration(url string, token string) error {
	dir := path.Join(os.Getenv("HOME"), ".peteq")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModePerm)
	}
	c := clientConfig{
		URL:   url,
		Token: token,
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(dir, "config"), data, os.ModePerm)
}
