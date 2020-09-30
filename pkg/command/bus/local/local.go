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
func (c *CommandBus) Execute(ctx context.Context, name string, arguments interface{}, done chan<- error) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	h, ok := c.Handlers[name]
	if !ok {
		done <- fmt.Errorf("Handler not found")
		return
	}
	c.Logger.Info("Calling command handler", "name", name)
	h.Handle(ctx, done, arguments)
}

// ExecuteAndWait execute command and waits for it to be completed
// return an an error if happen
func (c *CommandBus) ExecuteAndWait(ctx context.Context, name string, arguments interface{}) error {
	var err error
	wg := &sync.WaitGroup{}
	cn := make(chan error)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		for {
			select {
			case e := <-cn:
				if e != nil {
					err = e
				}
				wg.Done()
				return
			case _ = <-ctx.Done():
				wg.Done()
				return
			}
		}
	}(wg)
	c.Execute(ctx, name, arguments, cn)
	wg.Wait()
	return err

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
