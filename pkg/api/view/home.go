package view

import (
	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/task"
)

type (
	// HomeViewAPI for backlog view
	HomeViewAPI struct {
		TaskRepo    *task.Repo
		ListRepo    *list.Repo
		ProjectRepo *project.Repo
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
func (h *HomeViewAPI) Get(c *gin.Context) {
	lists, err := h.ListRepo.List(list.QueryOptions{})
	if err != nil {
		handleError(400, err, c)
		return
	}

	tasks, err := h.TaskRepo.List(task.ListOptions{})
	if err != nil {
		handleError(400, err, c)
		return
	}

	projects, err := h.ProjectRepo.List(project.QueryOptions{})
	if err != nil {
		handleError(400, err, c)
		return
	}
	homeLists := []homeList{}
	for _, l := range lists {
		homeTasks := []homeTask{}
		for _, t := range l.Tasks {
			for _, task := range tasks {
				if t == task.Metadata.ID {
					var p project.Project
					for _, proj := range projects {
						for _, tid := range proj.Tasks {
							if tid == task.Metadata.ID {
								p = proj
							}
						}
					}
					homeTasks = append(homeTasks, homeTask{task, p})
				}
			}
		}
		homeLists = append(homeLists, homeList{l, homeTasks})
	}
	c.JSON(200, homeView{Lists: homeLists})
}
