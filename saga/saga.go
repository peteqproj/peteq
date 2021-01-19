package saga

import (
	"context"

	automationDomain "github.com/peteqproj/peteq/domain/automation"
	listDomain "github.com/peteqproj/peteq/domain/list"
	projectDomain "github.com/peteqproj/peteq/domain/project"
	taskDomain "github.com/peteqproj/peteq/domain/task"
	triggerDomain "github.com/peteqproj/peteq/domain/trigger"
	triggerEventTypes "github.com/peteqproj/peteq/domain/trigger/event/types"
	userDomain "github.com/peteqproj/peteq/domain/user"
	userEventTypes "github.com/peteqproj/peteq/domain/user/event/types"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// Saga is long running process
	// saga may fail
	Saga interface {
		Run(context.Context) error
	}

	// EventHandler handle all the events that starts saga process
	EventHandler struct {
		ListRepo       *listDomain.Repo
		TaskRepo       *taskDomain.Repo
		AutomationRepo *automationDomain.Repo
		ProjectRepo    *projectDomain.Repo
		TriggerRepo    *triggerDomain.Repo
		UserRepo       *userDomain.Repo
		CommandBus     commandbus.CommandBus
	}
)

func (e *EventHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	logger.Info("Handling saga event", "event", ev.Metadata.Name, "id", ev.Metadata.ID)

	switch ev.Metadata.Name {
	case userEventTypes.UserRegistredEvent:
		{
			return (&registrator{
				Commandbus:  e.CommandBus,
				Logger:      logger,
				ListRepo:    e.ListRepo,
				IDGenerator: utils.NewGenerator(),
			}).Run(ctx)
		}
	case triggerEventTypes.TriggerTriggeredEvent:
		{
			tb, err := e.AutomationRepo.GetTriggerBindingByTriggerID(ev.Tenant.ID, ev.Metadata.AggregatorID)
			if err != nil {
				return err
			}
			a, err := e.AutomationRepo.Get(ev.Tenant.ID, tb.Spec.Automation)
			if err != nil {
				return err
			}
			switch a.Spec.Type {
			case "task-archiver":
				return newTaskArchiver(e.CommandBus, e.TaskRepo, e.ListRepo, logger, ev.Tenant.ID).Run(ctx)
			}
			logger.Info("Spec does not match to any known saga process", "type", a.Spec.Type)
		}
	}

	logger.Info("Event does not match to any known saga process", "event", ev.Metadata.Name)
	return nil

}
func (e *EventHandler) Name() string {
	return "saga_event_handler"
}
func newTaskArchiver(cb commandbus.CommandBus, taskRepo *taskDomain.Repo, listRepo *listDomain.Repo, lgr logger.Logger, user string) Saga {
	return &archiver{
		Commandbus: cb,
		TaskRepo:   taskRepo,
		ListRepo:   listRepo,
		Logger:     lgr,
		User:       user,
	}
}
