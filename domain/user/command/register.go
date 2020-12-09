package command

import (
	"context"
	"fmt"
	"time"

	automationCommand "github.com/peteqproj/peteq/domain/automation/command"
	listCommand "github.com/peteqproj/peteq/domain/list/command"
	triggerCommand "github.com/peteqproj/peteq/domain/trigger/command"
	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/domain/user/event/handler"
	"github.com/peteqproj/peteq/domain/user/event/types"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// RegisterCommand to create task
	RegisterCommand struct {
		Eventbus    bus.Eventbus
		Repo        *user.Repo
		Commandbus  commandbus.CommandBus
		IDGenerator utils.IDGenerator
	}

	// RegisterCommandOptions to create new user
	RegisterCommandOptions struct {
		UserID       string
		Email        string
		PasswordHash string
	}
)

// Handle runs RegisterCommand to create new user
func (r *RegisterCommand) Handle(ctx context.Context, arguments interface{}) error {
	opt, ok := arguments.(RegisterCommandOptions)
	if !ok {
		return fmt.Errorf("Failed to convert arguments to User")
	}
	usr, err := r.Repo.GetByEmail(opt.Email)
	if err != nil {
		if err.Error() != "User not found" {
			return err
		}
	}
	if usr != nil {
		return fmt.Errorf("Email already registred")
	}
	_, err = r.Eventbus.Publish(ctx, event.Event{
		Tenant: tenant.Tenant{
			ID:   opt.UserID,
			Type: tenant.User.String(),
		},
		Metadata: event.Metadata{
			Name:           types.UserRegistredEvent,
			CreatedAt:      time.Now(),
			AggregatorRoot: "user",
			AggregatorID:   opt.UserID,
		},
		Spec: handler.RegisteredSpec{
			Email:        opt.Email,
			ID:           opt.UserID,
			PasswordHash: opt.PasswordHash,
		},
	})

	basicLists := []string{"Upcoming", "Today", "Done"}
	ectx := tenant.ContextWithUser(ctx, user.User{
		Metadata: user.Metadata{
			Email: opt.Email,
			ID:    opt.UserID,
		},
	})
	for i, l := range basicLists {
		id, err := r.IDGenerator.GenerateV4()
		if err != nil {
			return err
		}
		if err := r.Commandbus.Execute(ectx, "list.create", listCommand.CreateCommandOptions{
			Name:  l,
			ID:    id,
			Index: i,
		}); err != nil {
			return err
		}
	}

	tid, err := r.IDGenerator.GenerateV4()
	if err != nil {
		return err
	}
	if err := r.Commandbus.Execute(ectx, "trigger.create", triggerCommand.TriggerCreateCommandOptions{
		ID:          tid,
		Name:        "Task Archiver",
		Description: "Runs at 00:00 every day",
		Cron:        utils.PtrString("0 00 * * 0-4"), // “At 00:00 on every day-of-week from Sunday through Thursday.”
	}); err != nil {
		return err
	}

	tid2, err := r.IDGenerator.GenerateV4()
	if err != nil {
		return err
	}
	if err := r.Commandbus.Execute(ectx, "automation.create", automationCommand.AutomationCreateCommandOptions{
		ID:              tid2,
		Name:            "Task Archiver",
		Description:     "Archive tasks in Done list",
		Type:            "task-archiver",
		JSONInputSchema: "",
	}); err != nil {
		return err
	}

	tid3, err := r.IDGenerator.GenerateV4()
	if err != nil {
		return err
	}
	if err := r.Commandbus.Execute(ectx, "automation.bindTrigger", automationCommand.TriggerBindingCreateCommandOptions{
		ID:         tid3,
		Name:       fmt.Sprintf("Bind Trigger \"%s\" to Automation \"%s\" ", "Task Archiver", "Task Archiver"),
		Automation: tid2,
		Trigger:    tid,
	}); err != nil {
		return err
	}
	return err
}
