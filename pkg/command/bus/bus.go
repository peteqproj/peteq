package bus

import (
	"context"
	"sync"

	"github.com/peteqproj/peteq/pkg/command/bus/local"
	"github.com/peteqproj/peteq/pkg/command/handler"
)

type (
	// CommandBus used to pass commands on bus
	CommandBus interface {
		Execute(ctx context.Context, name string, arguments interface{}, done chan<- error)
		ExecuteAndWait(ctx context.Context, name string, arguments interface{}) error
		RegisterHandler(name string, ch handler.CommandHandler) error
	}

	// Options to build commandbus
	Options struct {
		Type string
	}
)

// New builds commandbus from options
func New(options Options) CommandBus {
	if options.Type == "local" {
		return &local.CommandBus{
			Handlers: map[string]handler.CommandHandler{},
			Lock:     &sync.Mutex{},
		}
	}
	return nil
}
