package home

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/list"
	listCommand "github.com/peteqproj/peteq/domain/list/command"
	"github.com/peteqproj/peteq/domain/project"
	projectCommand "github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/logger"
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
	view, err := h.DAL.load(c.Request.Context(), u.Metadata.ID)
	if err != nil {
		handleError(400, err, c)
		return
	}
	c.JSON(200, view)
}
func (h *ViewAPI) EventHandlers() map[string]handler.EventHandler {
	return map[string]handler.EventHandler{
		"user.registred":     h,
		"list.created":       h,
		"list.task-moved":    h,
		"task.updated":       h,
		"task.deleted":       h,
		"project.task-added": h,
	}
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}

func (h *ViewAPI) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	if ev.Metadata.Name == "user.registred" {
		return h.handlerUserRegistration(ctx, ev, logger)
	}
	current, err := h.DAL.load(ctx, ev.Tenant.ID)
	if err != nil {
		return err
	}
	updated, err := h.handlerUpdateEvent(ctx, ev, current, logger)
	if err != nil {
		return err
	}
	return h.DAL.update(ctx, ev.Tenant.ID, updated)
}
func (h *ViewAPI) Name() string {
	return "home_view"
}

func (h *ViewAPI) handlerUpdateEvent(ctx context.Context, ev event.Event, view homeView, logger logger.Logger) (homeView, error) {
	switch ev.Metadata.Name {
	case "list.created":
		{
			return h.handlerListCreated(ctx, ev, view, logger)
		}
	case "list.task-moved":
		{
			return h.handlerTaskAddedToList(ctx, ev, view, logger)
		}
	case "task.updated":
		{
			return h.handlerTaskUpdated(ctx, ev, view, logger)
		}
	case "task.deleted":
		{
			return h.handlerTaskDeleted(ctx, ev, view, logger)
		}
	case "project.task-added":
		{
			return h.handlerTaskAddedToProject(ctx, ev, view, logger)
		}
	}
	return view, fmt.Errorf("Event handler not found")
}
func (h *ViewAPI) handlerUserRegistration(ctx context.Context, ev event.Event, logger logger.Logger) error {
	v := homeView{
		Lists: []homeList{},
	}
	return h.DAL.create(ctx, ev.Tenant.ID, v)
}
func (h *ViewAPI) handlerListCreated(ctx context.Context, ev event.Event, view homeView, logger logger.Logger) (homeView, error) {
	opt := listCommand.CreateCommandOptions{}
	if err := ev.UnmarshalSpecInto(&opt); err != nil {
		return view, err
	}
	found := false
	for _, l := range view.Lists {
		if l.Metadata.ID == opt.ID {
			found = true
		}
	}
	if found {
		logger.Info("List already added to view", "list", opt.ID)
		return view, nil
	}

	view.Lists = append(view.Lists, homeList{
		List: list.List{
			Tasks: []string{},
			Metadata: list.Metadata{
				ID:    opt.ID,
				Name:  opt.Name,
				Index: opt.Index,
			},
			Tenant: ev.Tenant,
		},
		Tasks: []homeTask{},
	})
	return view, nil
}
func (h *ViewAPI) handlerTaskAddedToList(ctx context.Context, ev event.Event, view homeView, logger logger.Logger) (homeView, error) {
	opt := listCommand.MoveTaskArguments{}
	if err := ev.UnmarshalSpecInto(&opt); err != nil {
		return view, err
	}
	task, err := h.TaskRepo.Get(ev.Tenant.ID, opt.TaskID)
	if err != nil {
		return view, err
	}
	sourceIndex := -1
	destinationIndex := -1
	for i, l := range view.Lists {
		if opt.Source != "" && l.Metadata.ID == opt.Source {
			sourceIndex = i
			continue
		}

		if opt.Destination != "" && l.Metadata.ID == opt.Destination {
			destinationIndex = i
			continue
		}
	}

	// search if there is reference for task in any project
	projects, err := h.ProjectRepo.List(project.QueryOptions{
		UserID: ev.Tenant.ID,
	})
	if err != nil {
		return view, err
	}

	projectIndex := -1
	taskInProjectIndex := -1
	for i, p := range projects {
		if taskInProjectIndex != -1 {
			break
		}
		for j, t := range p.Tasks {
			if t == opt.TaskID {
				projectIndex = i
				taskInProjectIndex = j
				break
			}
		}
	}
	taskProject := project.Project{}
	if taskInProjectIndex != -1 {
		taskProject = projects[projectIndex]
	}

	// If source found, remove task from source
	if sourceIndex != -1 {
		for i, tid := range view.Lists[sourceIndex].Tasks {
			if tid.Task.Metadata.ID == opt.TaskID {
				view.Lists[sourceIndex].Tasks = remove(view.Lists[sourceIndex].Tasks, i)
				break
			}
		}
	}

	// If destination found add it to destination
	if destinationIndex != -1 {
		view.Lists[destinationIndex].Tasks = append(view.Lists[destinationIndex].Tasks, homeTask{
			Task:    task,
			Project: taskProject,
		})
	}
	return view, nil
}
func (h *ViewAPI) handlerTaskUpdated(ctx context.Context, ev event.Event, view homeView, logger logger.Logger) (homeView, error) {
	task := task.Task{}
	err := ev.UnmarshalSpecInto(&task)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}
	listIndex, taskIndex := findTaskInView(view, task.Metadata.ID)
	if taskIndex == -1 {
		// task not in lists, no action to do
		return view, nil
	}
	view.Lists[listIndex].Tasks[taskIndex].Task = task
	return view, nil
}
func (h *ViewAPI) handlerTaskDeleted(ctx context.Context, ev event.Event, view homeView, logger logger.Logger) (homeView, error) {
	task := task.Task{}
	err := ev.UnmarshalSpecInto(&task)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}

	listIndex, taskIndex := findTaskInView(view, task.Metadata.ID)
	if taskIndex == -1 {
		// task not in lists
		return view, nil
	}
	view.Lists[listIndex].Tasks = append(view.Lists[listIndex].Tasks[:taskIndex], view.Lists[listIndex].Tasks[taskIndex+1:]...)
	return view, nil
}
func (h *ViewAPI) handlerTaskAddedToProject(ctx context.Context, ev event.Event, view homeView, logger logger.Logger) (homeView, error) {
	opt := projectCommand.AddTasksCommandOptions{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to AddTasksCommandOptions object: %v", err)
	}
	newProject, err := h.ProjectRepo.Get(ev.Tenant.ID, opt.Project)
	if err != nil {
		return view, err
	}
	listIndex, taskIndex := findTaskInView(view, opt.TaskID)
	if taskIndex == -1 {
		// task not found in lists, not an error
		return view, nil
	}
	view.Lists[listIndex].Tasks[taskIndex].Project = newProject
	return view, nil
}

func remove(slice []homeTask, s int) []homeTask {
	return append(slice[:s], slice[s+1:]...)
}
func findTaskInView(view homeView, id string) (int, int) {
	listIndex := -1
	taskIndex := -1
	for i, l := range view.Lists {
		for j, t := range l.Tasks {
			if t.Metadata.ID == id {
				listIndex = i
				taskIndex = j
				break
			}
		}
		if taskIndex != -1 {
			break
		}
	}
	return listIndex, taskIndex
}
