package view

import (
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/task"
)

type (
	// BacklogViewAPI for backlog view
	BacklogViewAPI struct {
		TaskRepo    *task.Repo
		ListRepo    *list.Repo
		ProjectRepo *project.Repo
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
func (b *BacklogViewAPI) Get(c *gin.Context) {
	listSet, err := b.ListRepo.List(list.QueryOptions{})
	if err != nil {
		handleError(400, err, c)
		return
	}
	projectSet, err := b.ProjectRepo.List(project.QueryOptions{})
	if err != nil {
		handleError(400, err, c)
		return
	}

	taskSet, err := b.TaskRepo.List(task.ListOptions{})
	if err != nil {
		handleError(400, err, c)
		return
	}

	backlogTasks := []backlogTask{}
	lists := []backlogTaskList{}
	for _, l := range listSet {
		lists = append(lists, backlogTaskList{
			ID:   l.Metadata.ID,
			Name: l.Metadata.Name,
		})
	}
	projects := []backlogTaskProject{}
	for _, p := range projectSet {
		projects = append(projects, backlogTaskProject{
			ID:   p.Metadata.ID,
			Name: p.Metadata.Name,
		})
	}
	sort.SliceStable(taskSet, func(i, j int) bool {
		t1 := taskSet[i]
		return !t1.Status.Completed
	})
	for _, task := range taskSet {
		backlogList := backlogTaskList{
			ID:   "-1",
			Name: "=====",
		}
		for _, list := range listSet {
			for _, id := range list.Tasks {
				if id == task.Metadata.ID {
					backlogList.ID = list.Metadata.ID
					backlogList.Name = list.Metadata.Name
				}
			}
		}

		backlogProject := backlogTaskProject{
			ID:   "-1",
			Name: "=====",
		}
		for _, proj := range projectSet {
			for _, id := range proj.Tasks {
				if id == task.Metadata.ID {
					backlogProject.ID = proj.Metadata.ID
					backlogProject.Name = proj.Metadata.Name
				}
			}
		}
		backlogTasks = append(backlogTasks, backlogTask{task, backlogList, backlogProject})
	}

	c.JSON(200, backlogView{backlogTasks, lists, projects})
}
