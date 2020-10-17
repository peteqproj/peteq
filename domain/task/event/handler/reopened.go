package handler

import (
	"context"
	"fmt"

	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// ReopenedHandler to handle task.created event
	ReopenedHandler struct {
		Repo *task.Repo
	}
)

// Handle will handle the event the process it
func (c *ReopenedHandler) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	task, err := c.Repo.Get(ev.Tenant.ID, ev.Metadata.AggregatorID)
	if err != nil {
		return fmt.Errorf("Failed to get task %s: %v", ev.Metadata.AggregatorID, err)
	}
	task.Status.Completed = false
	return c.Repo.Update(ev.Tenant.ID, task)
}

func (c *ReopenedHandler) Name() string {
	return "domain_ReopenedHandler"
}
