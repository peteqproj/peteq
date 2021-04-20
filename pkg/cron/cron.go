package cron

import (
	"context"
	"time"

	"github.com/peteqproj/peteq/domain/trigger/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	eventbus "github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
	crontab "github.com/robfig/cron/v3"
)

type (
	// Cron run tasks continouosly
	Cron interface {
		AddFunc(trigger string, cronExp string)
		Start()
		Stop()
	}

	// Options to create new Cron
	Options struct {
		EventPublisher eventbus.EventPublisher
		Logger         logger.Logger
		UserID         string
	}

	cr struct {
		cron           *crontab.Cron
		eventPublisher eventbus.EventPublisher
		logger         logger.Logger
		userID         string
	}
)

// New creates new Cron from Options
func New(opt Options) Cron {
	return &cr{
		cron:           crontab.New(),
		eventPublisher: opt.EventPublisher,
		logger:         opt.Logger,
		userID:         opt.UserID,
	}
}

func (c *cr) AddFunc(trigger string, cronExp string) {
	if _, err := c.cron.AddFunc(cronExp, c.handleCronTick(trigger)); err != nil {
		c.logger.Info("Failed to add cron func", "error", err.Error())
		return
	}
	c.logger.Info("Cron func added")
}

func (c *cr) Start() {
	c.cron.Run()
}
func (c *cr) Stop() {
	ctx := c.cron.Stop()
	select {
	case _ = <-ctx.Done():
		{
			c.logger.Info("Cron stopped")
		}
	case _ = <-time.After(time.Second * 30):
		{
			c.logger.Info("Cron wasnt not stopped normally, force after 30 seconds")
			return
		}
	}
}

func (c *cr) handleCronTick(trigger string) func() {
	return func() {
		c.logger.Info("Triggered", "user", c.userID, "trigger", trigger)
		_, err := c.eventPublisher.Publish(context.Background(), event.Event{
			Tenant: tenant.Tenant{
				ID:   c.userID,
				Type: tenant.User.String(),
			},
			Metadata: event.Metadata{
				AggregatorRoot: "trigger",
				AggregatorID:   trigger,
				Name:           types.TriggerTriggeredEvent,
				CreatedAt:      time.Now(),
			},
		})
		if err != nil {
			c.logger.Info("Failed to publish event", "error", err)
		}
	}
}
