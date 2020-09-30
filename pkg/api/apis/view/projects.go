package view

import (
	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// ProjectsViewAPI for projects view
	ProjectsViewAPI struct {
		TaskRepo    *task.Repo
		ProjectRepo *project.Repo
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
func (b *ProjectsViewAPI) Get(c *gin.Context) {
	u := tenant.UserFromContext(c.Request.Context())
	projects, err := b.ProjectRepo.List(project.QueryOptions{UserID: u.Metadata.ID})
	if err != nil {
		handleError(400, err, c)
		return
	}

	taskSet, err := b.TaskRepo.List(task.ListOptions{UserID: u.Metadata.ID})
	if err != nil {
		handleError(400, err, c)
		return
	}
	populatedProjects := []populatedProject{}
	for _, proj := range projects {
		project := populatedProject{
			Project: proj,
			Tasks:   []task.Task{},
		}
		for _, t := range proj.Tasks {
			for _, task := range taskSet {
				if t == task.Metadata.ID {
					project.Tasks = append(project.Tasks, task)
				}
			}
		}
		populatedProjects = append(populatedProjects, project)
	}

	c.JSON(200, projectsView{Projects: populatedProjects})
}
