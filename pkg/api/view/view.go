package view

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/api/view/backlog"
	"github.com/peteqproj/peteq/pkg/api/view/home"
	proj "github.com/peteqproj/peteq/pkg/api/view/project"
	"github.com/peteqproj/peteq/pkg/api/view/projects"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// View can only retrieve data
	View interface {
		Get(c *gin.Context)
		EventHandlers() map[string]handler.EventHandler
	}

	// Options to build view
	Options struct {
		Type        string
		TaskRepo    *task.Repo
		ListRepo    *list.Repo
		ProjectRepo *project.Repo
		Logger      logger.Logger
		DB          *sql.DB
	}
)

// NewView builds view from options
func NewView(opt Options) View {
	switch opt.Type {
	case "home":
		return &home.ViewAPI{
			TaskRepo:    opt.TaskRepo,
			ListRepo:    opt.ListRepo,
			ProjectRepo: opt.ProjectRepo,
			DAL: &home.DAL{
				DB: opt.DB,
			},
		}
	case "backlog":
		return &backlog.ViewAPI{
			TaskRepo:    opt.TaskRepo,
			ListRepo:    opt.ListRepo,
			ProjectRepo: opt.ProjectRepo,
			DAL: &backlog.DAL{
				DB: opt.DB,
			},
		}
	case "projects":
		return &projects.ViewAPI{
			TaskRepo:    opt.TaskRepo,
			ProjectRepo: opt.ProjectRepo,
			DAL: &projects.DAL{
				DB: opt.DB,
			},
		}
	case "project":
		return &proj.ViewAPI{
			TaskRepo:    opt.TaskRepo,
			ProjectRepo: opt.ProjectRepo,
			DAL: &proj.DAL{
				DB: opt.DB,
			},
		}
	}
	return nil
}
