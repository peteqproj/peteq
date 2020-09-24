package view

import (
	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/task"
)

type (
	// View can only retrieve data
	View interface {
		Get(c *gin.Context)
	}

	// Options to build view
	Options struct {
		Type        string
		TaskRepo    *task.Repo
		ListRepo    *list.Repo
		ProjectRepo *project.Repo
	}
)

// NewView builds view from options
func NewView(opt Options) View {
	switch opt.Type {
	case "home":
		return &HomeViewAPI{
			TaskRepo:    opt.TaskRepo,
			ListRepo:    opt.ListRepo,
			ProjectRepo: opt.ProjectRepo,
		}
	case "backlog":
		return &BacklogViewAPI{
			TaskRepo:    opt.TaskRepo,
			ListRepo:    opt.ListRepo,
			ProjectRepo: opt.ProjectRepo,
		}
	case "projects":
		return &ProjectsViewAPI{
			TaskRepo:    opt.TaskRepo,
			ProjectRepo: opt.ProjectRepo,
		}
	case "project":
		return &ProjectViewAPI{
			TaskRepo:    opt.TaskRepo,
			ProjectRepo: opt.ProjectRepo,
		}
	}
	return nil
}
