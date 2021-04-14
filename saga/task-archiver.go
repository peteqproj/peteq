package saga

import (
	"context"
	"fmt"

	listDomain "github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/task"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	archiver struct {
		Commandbus commandbus.CommandBus
		TaskRepo   *task.Repo
		ListRepo   *listDomain.Repo
		Logger     logger.Logger
		User       string // TODO: use from context tenant
	}
)

func (a *archiver) Run(ctx context.Context) error {
	a.Logger.Info("Running task archiver")
	// TODO: make sure the context is authenticated and remove a.User
	lists, err := a.ListRepo.ListByUserid(ctx, a.User)
	if err != nil {
		return fmt.Errorf("Failed to get lists: %w", err)
	}
	candidates := []string{}
	for _, l := range lists {
		if l.Metadata.Name != "Done" {
			continue
		}
		candidates = l.Spec.Tasks
	}

	for _, c := range candidates {
		t, err := a.TaskRepo.GetById(ctx, c)
		if err != nil {
			a.Logger.Info("Failed to request task", "id", c, "error", err.Error())
			continue
		}
		a.Logger.Info("Deleting task", "task", t.Metadata.ID)
		if err := a.Commandbus.Execute(ctx, "task.delete", t); err != nil {
			return err
		}
	}
	return nil
}
