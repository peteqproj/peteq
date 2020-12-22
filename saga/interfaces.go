package saga

import (
	"github.com/peteqproj/peteq/domain/list"
)

type (
	// ListRepo common interface for all sagas
	ListRepo interface {
		GetListByName(userID string, name string) (list.List, error)
	}
)
