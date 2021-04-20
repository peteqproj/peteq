package saga

import (
	"context"

	"github.com/peteqproj/peteq/domain/list"
)

type (
	// ListRepo common interface for all sagas
	ListRepo interface {
		GetById(ctx context.Context, id string) (*list.List, error)
	}
)
