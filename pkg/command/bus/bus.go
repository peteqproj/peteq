package bus

import (
	"context"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/peteqproj/peteq/pkg/command/bus/google"
	"github.com/peteqproj/peteq/pkg/command/bus/local"
	"github.com/peteqproj/peteq/pkg/command/handler"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// CommandBus used to pass commands on bus
	CommandBus interface {
		Execute(ctx context.Context, name string, arguments interface{}) error
		RegisterHandler(name string, ch handler.CommandHandler) error
		Start() error
	}

	// Options to build commandbus
	Options struct {
		Type                     string
		Logger                   logger.Logger
		ExtendContextFunc        func(context.Context, string) context.Context
		GooglePubSubClient       *pubsub.Client
		GooglePubSubTopic        string
		GooglePubSubSubscribtion string
	}
)

// New builds commandbus from options
func New(options Options) CommandBus {
	if options.Type == "local" {
		return &local.CommandBus{
			Handlers: map[string]handler.CommandHandler{},
			Lock:     &sync.Mutex{},
			Logger:   options.Logger,
		}
	}

	if options.Type == "google" {
		return &google.Bus{
			Handlers:          map[string]handler.CommandHandler{},
			Lock:              &sync.Mutex{},
			Logger:            options.Logger,
			Ps:                options.GooglePubSubClient,
			Topic:             options.GooglePubSubTopic,
			Subscribtion:      options.GooglePubSubSubscribtion,
			IDGenerator:       utils.NewGenerator(),
			ExtendContextFunc: options.ExtendContextFunc,
		}
	}
	return nil
}
