package local

import (
	"context"
	"fmt"
	"sync"

	"github.com/peteqproj/peteq/pkg/command/handler"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// CommandBus local, in memory bus
	CommandBus struct {
		Handlers map[string]handler.CommandHandler
		Lock     *sync.Mutex
		Logger   logger.Logger
	}
)

// Execute runs command handler
// error will be reported to the channel
func (c *CommandBus) Execute(ctx context.Context, name string, arguments interface{}) error {
	h, ok := c.Handlers[name]
	if !ok {
		return fmt.Errorf("Handler not found")
	}
	c.Logger.Info("Calling command handler", "name", name)
	return h.Handle(ctx, arguments)
}

// RegisterHandler registers new command handler
func (c *CommandBus) RegisterHandler(name string, ch handler.CommandHandler) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	_, ok := c.Handlers[name]
	if ok {
		return fmt.Errorf("Handler already exist")
	}
	c.Handlers[name] = ch
	return nil
}
