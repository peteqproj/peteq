package backlog

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

	backlogView struct {
		Tasks    []backlogTask        `json:"tasks"`
		Lists    []backlogTaskList    `json:"lists"`
		Projects []backlogTaskProject `json:"projects"`
	}

	backlogTask struct {
		task.Task
		List    backlogTaskList    `json:"list"`
		Project backlogTaskProject `json:"project"`
	}

	backlogTaskList struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}

	backlogTaskProject struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}
)

// Get build backlog view
func (b *ViewAPI) Get(c *gin.Context) {
	u := tenant.UserFromContext(c.Request.Context())
	view, err := b.DAL.Get(c.Request.Context(), u.Metadata.ID)
	if err != nil {
		handleError(400, err, c)
		return
	}

	c.JSON(200, view)
}

func (h *ViewAPI) EventHandlers() map[string]handler.EventHandler {
	return map[string]handler.EventHandler{
		"list.task-moved": listTaskMovedHandler{
			dal:      h.DAL,
			listRepo: h.ListRepo,
			taskRepo: h.TaskRepo,
		},
		"task.created": taskCreatedHandler{
			dal: h.DAL,
		},
		"task.updated": taskUpdateHandler{
			dal: h.DAL,
		},
		"task.completed": taskStatusChangedHandler{
			dal:      h.DAL,
			taskRepo: h.TaskRepo,
		},
		"task.reopened": taskStatusChangedHandler{
			dal:      h.DAL,
			taskRepo: h.TaskRepo,
		},
		"task.deleted": taskDeletedHandler{
			dal: h.DAL,
		},
		"user.registred": userRegistredHandler{
			dal: h.DAL,
		},
		"project.created": projectCreatedHandler{
			dal:         h.DAL,
			projectRepo: h.ProjectRepo,
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
