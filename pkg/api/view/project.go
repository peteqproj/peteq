package view

import (
	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/task"
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
	id := c.Param("id")
	proj, err := b.ProjectRepo.Get(id)
	if err != nil {
		handleError(400, err, c)
		return
	}

	taskSet, err := b.TaskRepo.List(task.ListOptions{})
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
