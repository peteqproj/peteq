package saga

import (
	"context"
	"fmt"

	automationCommand "github.com/peteqproj/peteq/domain/automation/command"
	listDomain "github.com/peteqproj/peteq/domain/list"
	listCommand "github.com/peteqproj/peteq/domain/list/command"
	triggerCommand "github.com/peteqproj/peteq/domain/trigger/command"
	userDomain "github.com/peteqproj/peteq/domain/user"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	registrator struct {
		Commandbus  commandbus.CommandBus
		ListRepo    *listDomain.Repo
		UserRepo    *userDomain.Repo
		Logger      logger.Logger
		IDGenerator utils.IDGenerator
	}
)

func (a *registrator) Run(ctx context.Context) error {
	a.Logger.Info("Running user registrator")
	user := tenant.UserFromContext(ctx)
	if user == nil {
		return fmt.Errorf("Failed to register user. Saga context does not include user")
	}
	if err := a.createBasicLists(ctx); err != nil {
		return fmt.Errorf("Failed to create basic lists: %w", err)
	}

	if err := a.createBasicTriggerAndAutomation(ctx); err != nil {
		return fmt.Errorf("Failed to create basic trigger or automation: %w", err)
	}
	return nil
}

func (a *registrator) createBasicLists(ctx context.Context) error {
	// TODO: check if those lists already exists
	basicLists := []string{"Upcoming", "Today", "Done"}
	for i, l := range basicLists {
		id, err := a.IDGenerator.GenerateV4()
		if err != nil {
			return err
		}
		if err := a.Commandbus.Execute(ctx, "list.create", listCommand.CreateCommandOptions{
			Name:  l,
			ID:    id,
			Index: i,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (a *registrator) createBasicTriggerAndAutomation(ctx context.Context) error {

	tid, err := a.IDGenerator.GenerateV4()
	if err != nil {
		return err
	}
	if err := a.Commandbus.Execute(ctx, "trigger.create", triggerCommand.TriggerCreateCommandOptions{
		ID:          tid,
		Name:        "Task Archiver",
		Description: "Runs at 00:00 every day",
		Cron:        utils.PtrString("0 00 * * 0-4"), // “At 00:00 on every day-of-week from Sunday through Thursday.”
	}); err != nil {
		return err
	}

	tid2, err := a.IDGenerator.GenerateV4()
	if err != nil {
		return err
	}
	if err := a.Commandbus.Execute(ctx, "automation.create", automationCommand.AutomationCreateCommandOptions{
		ID:              tid2,
		Name:            "Task Archiver",
		Description:     "Archive tasks in Done list",
		Type:            "task-archiver",
		JSONInputSchema: "",
	}); err != nil {
		return err
	}

	tid3, err := a.IDGenerator.GenerateV4()
	if err != nil {
		return err
	}
	if err := a.Commandbus.Execute(ctx, "automation.bindTrigger", automationCommand.TriggerBindingCreateCommandOptions{
		ID:         tid3,
		Name:       fmt.Sprintf("Bind Trigger \"%s\" to Automation \"%s\" ", "Task Archiver", "Task Archiver"),
		Automation: tid2,
		Trigger:    tid,
	}); err != nil {
		return err
	}
	return nil
}
