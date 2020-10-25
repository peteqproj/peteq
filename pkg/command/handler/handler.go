package handler

import "context"

type (
	// CommandHandler runs command request
	CommandHandler interface {
		Handle(ctx context.Context, arguments interface{}) error
	}
)
