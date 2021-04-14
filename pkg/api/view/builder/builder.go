package builder

import (
	"github.com/peteqproj/peteq/pkg/api/view"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/repo"

	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/api/view/views/backlog"
	"github.com/peteqproj/peteq/pkg/api/view/views/home"
	proj "github.com/peteqproj/peteq/pkg/api/view/views/project"
	"github.com/peteqproj/peteq/pkg/api/view/views/projects"
	"github.com/peteqproj/peteq/pkg/api/view/views/triggers"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	Builder struct {
		TaskRepo    *task.Repo
		ListRepo    *list.Repo
		ProjectRepo *repo.Repo
		Logger      logger.Logger
		DB          db.Database
	}

	Options struct {
		TaskRepo    *task.Repo
		ListRepo    *list.Repo
		ProjectRepo *repo.Repo
		Logger      logger.Logger
		DB          db.Database
	}
)

// New creates new View Builder
func New(opt *Options) *Builder {
	return &Builder{
		TaskRepo:    opt.TaskRepo,
		ListRepo:    opt.ListRepo,
		ProjectRepo: opt.ProjectRepo,
		Logger:      opt.Logger,
		DB:          opt.DB,
	}
}

// BuildViews build all views
func (b *Builder) BuildViews() map[string]view.View {
	views := map[string]view.View{}
	views["home"] = &home.ViewAPI{
		TaskRepo:    b.TaskRepo,
		ListRepo:    b.ListRepo,
		ProjectRepo: b.ProjectRepo,
		DAL: &home.DAL{
			DB: b.DB,
		},
	}
	views["backlog"] = &backlog.ViewAPI{
		TaskRepo:    b.TaskRepo,
		ListRepo:    b.ListRepo,
		ProjectRepo: b.ProjectRepo,
		DAL: &backlog.DAL{
			DB: b.DB,
		},
	}

	views["projects"] = &projects.ViewAPI{
		TaskRepo:    b.TaskRepo,
		ProjectRepo: b.ProjectRepo,
		DAL: &projects.DAL{
			DB: b.DB,
		},
	}

	views["projects/:id"] = &proj.ViewAPI{
		TaskRepo:    b.TaskRepo,
		ProjectRepo: b.ProjectRepo,
		DAL: &proj.DAL{
			DB: b.DB,
		},
	}

	views["triggers"] = &triggers.ViewAPI{
		DAL: &triggers.DAL{
			DB: b.DB,
		},
	}
	return views
}
