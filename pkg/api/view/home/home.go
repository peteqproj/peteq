package home

import (
	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// ViewAPI for backlog view
	ViewAPI struct {
		TaskRepo    *task.Repo
		ListRepo    *list.Repo
		ProjectRepo *project.Repo
		DAL         *DAL
	}

	homeView struct {
		Lists []homeList `json:"lists"`
	}

	homeList struct {
		list.List
		Tasks []homeTask `json:"tasks"`
	}

	homeTask struct {
		task.Task
		Project project.Project `json:"project"`
	}
)

// Get builds home view
func (h *ViewAPI) Get(c *gin.Context) {
	u := tenant.UserFromContext(c.Request.Context())
	view, err := h.DAL.Get(c.Request.Context(), u.Metadata.ID)
	if err != nil {
		handleError(400, err, c)
		return
	}
	c.JSON(200, view)
}

func (h *ViewAPI) EventHandlers() map[string]handler.EventHandler {
	return map[string]handler.EventHandler{
		"list.created": listCreatedHandler{
			dal: h.DAL,
		},
		"list.task-moved": listTaskMovedHandler{
			dal:         h.DAL,
			taskRepo:    h.TaskRepo,
			projectRepo: h.ProjectRepo,
		},
		"task.updated": taskUpdateHandler{
			dal: h.DAL,
		},
		"task.deleted": taskDeletedHandler{
			dal: h.DAL,
		},
		"user.registred": userRegistredHandler{
			dal: h.DAL,
		},
		"project.task-added": projectTaskAddedHandler{
			dal:         h.DAL,
			projectRepo: h.ProjectRepo,
			taskRepo:    h.TaskRepo,
		},
	}
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}
