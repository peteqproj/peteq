package view

import (
	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// ProjectViewAPI for single project view
	ProjectViewAPI struct {
		TaskRepo    *task.Repo
		ProjectRepo *project.Repo
	}

	projectView struct {
		project.Project
		Tasks []task.Task `json:"tasks"`
	}
)

// Get build project view
func (b *ProjectViewAPI) Get(c *gin.Context) {
	u := tenant.UserFromContext(c.Request.Context())
	id := c.Param("id")
	proj, err := b.ProjectRepo.Get(u.Metadata.ID, id)
	if err != nil {
		handleError(400, err, c)
		return
	}

	taskSet, err := b.TaskRepo.List(task.ListOptions{UserID: u.Metadata.ID})
	if err != nil {
		handleError(400, err, c)
		return
	}
	tasks := []task.Task{}
	for _, t := range proj.Tasks {
		for _, task := range taskSet {
			if t == task.Metadata.ID {
				tasks = append(tasks, task)
			}
		}
	}

	c.JSON(200, projectView{
		Project: proj,
		Tasks:   tasks,
	})
}
