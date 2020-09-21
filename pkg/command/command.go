package command

import "context"

type (
	// Command spesifies functions for command
	Command interface {
		Run(ctx context.Context, query interface{}) error
	}
)
