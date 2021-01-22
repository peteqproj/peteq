package local

import (
	"context"
	"encoding/json"
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

// Start local do nothing
func (c *CommandBus) Start() error {
	return nil
}

// Execute runs command handler
func (c *CommandBus) Execute(ctx context.Context, name string, arguments interface{}) error {
	h, ok := c.Handlers[name]
	if !ok {
		return fmt.Errorf("Handler not found")
	}
	c.Logger.Info("Calling command handler", "name", name)
	data, err := json.Marshal(arguments)
	if err != nil {
		return err
	}
	return h.Handle(ctx, data)
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
