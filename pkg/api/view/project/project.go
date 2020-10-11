package project

import (
	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// ViewAPI for single project view
	ViewAPI struct {
		TaskRepo    *task.Repo
		ProjectRepo *project.Repo
		DAL         *DAL
	}

	projectView struct {
		project.Project
		Tasks []task.Task `json:"tasks"`
	}
)

// Get build project view
func (b *ViewAPI) Get(c *gin.Context) {
	u := tenant.UserFromContext(c.Request.Context())
	id := c.Param("id")
	view, err := b.DAL.Get(c.Request.Context(), u.Metadata.ID, id)
	if err != nil {
		handleError(400, err, c)
		return
	}
	c.JSON(200, view)
}

func (h *ViewAPI) EventHandlers() map[string]handler.EventHandler {
	return map[string]handler.EventHandler{
		"task.deleted": taskDeletedHandler{
			dal:         h.DAL,
			projectRepo: h.ProjectRepo,
			taskRepo:    h.TaskRepo,
		},
		"project.task-added": projectTaskAddedHandler{
			dal:         h.DAL,
			projectRepo: h.ProjectRepo,
			taskRepo:    h.TaskRepo,
		},
		"project.created": projectCreatedHandler{
			dal: h.DAL,
		},
		"task.completed": taskStatusChangedHandler{
			dal:         h.DAL,
			taskRepo:    h.TaskRepo,
			projectRepo: h.ProjectRepo,
		},
		"task.reopened": taskStatusChangedHandler{
			dal:         h.DAL,
			taskRepo:    h.TaskRepo,
			projectRepo: h.ProjectRepo,
		},
	}
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}
