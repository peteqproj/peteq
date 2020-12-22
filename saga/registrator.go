package saga

import (
	"context"
	"fmt"

	automationCommand "github.com/peteqproj/peteq/domain/automation/command"
	listCommand "github.com/peteqproj/peteq/domain/list/command"
	triggerCommand "github.com/peteqproj/peteq/domain/trigger/command"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/errors"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	registrator struct {
		Commandbus  commandbus.CommandBus
		Logger      logger.Logger
		IDGenerator utils.IDGenerator
		ListRepo    ListRepo
	}
)

func (a *registrator) Run(ctx context.Context) error {
	a.Logger.Info("Running user registrator")
	user := tenant.UserFromContext(ctx)
	if user == nil {
		return fmt.Errorf("Authentication Error: user not found in context")
	}
	if err := a.createBasicLists(ctx); err != nil {
		return err
	}

	if err := a.createBasicTriggerAndAutomation(ctx); err != nil {
		return err
	}
	return nil
}

func (a *registrator) createBasicLists(ctx context.Context) error {
	basicLists := []string{"Upcoming", "Today", "Done"}
	for i, l := range basicLists {
		list, err := a.ListRepo.GetListByName(tenant.UserFromContext(ctx).Metadata.ID, l)

		var e errors.NotFoundError
		if err != nil && !errors.As(err, &e) {
			return err
		}

		id, err := a.IDGenerator.GenerateV4()
		if err != nil {
			return err
		}
		if err := a.Commandbus.Execute(ctx, "list.create", listCommand.CreateCommandOptions{
			Name:  list.Metadata.Name,
			ID:    id,
			Index: i,
		}); err != nil {
			return fmt.Errorf("Failed to create list %s: %w", l, err)
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
		return fmt.Errorf("Failed to create trigger Task Archiver: %w", err)
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
		return fmt.Errorf("Failed to create automation Task Archiver: %w", err)
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
		return fmt.Errorf("Failed to automation-trigger-binding for Task Archiver trigger: %w", err)
	}
	return nil
}
