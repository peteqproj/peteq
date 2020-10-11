package projects

import (
	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// ViewAPI for projects view
	ViewAPI struct {
		TaskRepo    *task.Repo
		ProjectRepo *project.Repo
		DAL         *DAL
	}

	projectsView struct {
		Projects []populatedProject `json:"projects"`
	}

	populatedProject struct {
		project.Project
		Tasks []task.Task `json:"tasks"`
	}
)

// Get build projects view
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
		"task.deleted": taskDeletedHandler{
			dal: h.DAL,
		},
		"project.task-added": projectTaskAddedHandler{
			dal:         h.DAL,
			projectRepo: h.ProjectRepo,
			taskRepo:    h.TaskRepo,
		},
		"project.created": projectCreatedHandler{
			dal: h.DAL,
		},
		"user.registred": userRegistredHandler{
			dal: h.DAL,
		},
	}
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}
