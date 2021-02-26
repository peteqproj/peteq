package saga

import (
	"context"
	"fmt"

	listDomain "github.com/peteqproj/peteq/domain/list"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/repo"
)

type (
	archiver struct {
		Commandbus commandbus.CommandBus
		TaskRepo   *repo.Repo
		ListRepo   *listDomain.Repo
		Logger     logger.Logger
		User       string // TODO: use from context tenant
	}
)

func (a *archiver) Run(ctx context.Context) error {
	a.Logger.Info("Running task archiver")
	lists, err := a.ListRepo.List(listDomain.QueryOptions{
		UserID: a.User,
	})
	if err != nil {
		return fmt.Errorf("Failed to get lists: %w", err)
	}
	candidates := []string{}
	for _, l := range lists {
		if l.Metadata.Name != "Done" {
			continue
		}
		candidates = l.Tasks
	}

	for _, c := range candidates {
		t, err := a.TaskRepo.Get(ctx, repo.GetOptions{ID: c})
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
