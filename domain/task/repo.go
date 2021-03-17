package task

import (
	"context"
	"errors"

	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/logger"
)

const db_name = "repo_task"

var errNotFound = errors.New("Resource not found")

type (
	Repo struct {
		DB     db.Database
		Logger logger.Logger
	}

	ListOptions struct{}
	GetOptions  struct {
		ID    string
		Query string
	}
)

func (r *Repo) List(ctx context.Context, options ListOptions) ([]*Task, error) {
	return nil, nil
}

func (r *Repo) Get(ctx context.Context, options GetOptions) (*Task, error) {
	return nil, nil
}

func (r *Repo) Create(ctx context.Context, tresource *Task) error {
	return nil
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *Repo) Update(ctx context.Context, resource *Task) error {
	return nil
}
