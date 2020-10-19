package backlog

import (
	"context"
	"encoding/json"
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
		"list.task-moved":    h,
		"task.created":       h,
		"task.updated":       h,
		"task.completed":     h,
		"task.reopened":      h,
		"task.deleted":       h,
		"user.registred":     h,
		"project.created":    h,
		"project.task-added": h,
	}
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
	return "backlog_view"
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}

func (h *ViewAPI) handlerUserRegistration(ctx context.Context, ev event.Event, logger logger.Logger) error {
	view := backlogView{
		Tasks:    make([]backlogTask, 0),
		Lists:    make([]backlogTaskList, 0),
		Projects: make([]backlogTaskProject, 0),
	}
	bytes, err := json.Marshal(view)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &view); err != nil {
		return err
	}
	return h.DAL.create(ctx, ev.Tenant.ID, view)
}

func (h *ViewAPI) handlerUpdateEvent(ctx context.Context, ev event.Event, view backlogView, logger logger.Logger) (backlogView, error) {
	switch ev.Metadata.Name {
	case "list.task-moved":
		{
			return h.handleTaskMovedToList(ctx, ev, view, logger)
		}
	case "task.created":
		{
			return h.handleTaskCreated(ctx, ev, view, logger)
		}
	case "task.updated":
		{
			return h.handleTaskUpdated(ctx, ev, view, logger)
		}
	case "task.completed":
		{
			return h.handleTaskStatusChanged(ctx, ev, view, logger)
		}
	case "task.reopened":
		{
			return h.handleTaskStatusChanged(ctx, ev, view, logger)
		}
	case "task.deleted":
		{
			return h.handleTaskDeleted(ctx, ev, view, logger)
		}
	case "project.created":
		{
			return h.handleProjectCreated(ctx, ev, view, logger)
		}
	case "project.task-added":
		{
			return h.handleTaskAddedToProject(ctx, ev, view, logger)
		}
	}
	return view, fmt.Errorf("Event handler not found")
}

func (h *ViewAPI) handleTaskMovedToList(ctx context.Context, ev event.Event, view backlogView, logger logger.Logger) (backlogView, error) {
	opt := listCommand.MoveTaskArguments{}
	if err := ev.UnmarshalSpecInto(&opt); err != nil {
		return view, err
	}
	newList := list.List{}
	if opt.Destination != "" {
		list, err := h.ListRepo.Get(ev.Tenant.ID, opt.Destination)
		if err != nil {
			return view, err
		}
		logger.Info("Destination is set", "id", opt.Destination, "name", list.Metadata.Name)
		newList = list
	}
	taskIndex := -1
	for i, t := range view.Tasks {
		if t.Task.Metadata.ID == opt.TaskID {
			logger.Info("Task found in view", "id", opt.TaskID, "index", i)
			taskIndex = i
		}
	}
	if taskIndex == -1 {
		return view, fmt.Errorf("Task not found")
	}
	view.Tasks[taskIndex].List = backlogTaskList{
		ID:   newList.Metadata.ID,
		Name: newList.Metadata.Name,
	}
	return view, nil
}
func (h *ViewAPI) handleTaskCreated(ctx context.Context, ev event.Event, view backlogView, logger logger.Logger) (backlogView, error) {
	task := task.Task{}
	err := ev.UnmarshalSpecInto(&task)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}
	view.Tasks = append(view.Tasks, backlogTask{
		Task: task,
	})
	return view, nil
}
func (h *ViewAPI) handleTaskUpdated(ctx context.Context, ev event.Event, view backlogView, logger logger.Logger) (backlogView, error) {
	task := task.Task{}
	err := ev.UnmarshalSpecInto(&task)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}
	index := findTask(view, task.Metadata.ID)
	if index == -1 {
		return view, fmt.Errorf("Task not found")
	}
	view.Tasks[index].Task = task
	return view, nil
}
func (h *ViewAPI) handleTaskStatusChanged(ctx context.Context, ev event.Event, view backlogView, logger logger.Logger) (backlogView, error) {
	task, err := h.TaskRepo.Get(ev.Tenant.ID, ev.Metadata.AggregatorID)
	if err != nil {
		return view, err
	}
	index := findTask(view, task.Metadata.ID)
	if index == -1 {
		return view, fmt.Errorf("Task not found")
	}
	view.Tasks[index].Task = task
	return view, nil
}
func (h *ViewAPI) handleTaskDeleted(ctx context.Context, ev event.Event, view backlogView, logger logger.Logger) (backlogView, error) {
	task := task.Task{}
	err := ev.UnmarshalSpecInto(&task)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}

	taskIndex := -1
	for i, t := range view.Tasks {
		if t.Metadata.ID == task.Metadata.ID {
			taskIndex = i
			break
		}
	}
	if taskIndex == -1 {
		return view, fmt.Errorf("Task not found")
	}
	view.Tasks = remove(view.Tasks, taskIndex)
	return view, nil
}
func (h *ViewAPI) handleTaskAddedToProject(ctx context.Context, ev event.Event, view backlogView, logger logger.Logger) (backlogView, error) {
	opt := projectCommand.AddTasksCommandOptions{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to task object: %v", err)
	}
	newProject := backlogTaskProject{}
	if opt.Project != "" {
		prj, err := h.ProjectRepo.Get(ev.Tenant.ID, opt.Project)
		if err != nil {
			return view, err
		}
		newProject = backlogTaskProject{
			ID:   prj.Metadata.ID,
			Name: prj.Metadata.Name,
		}
	}
	taskIndex := -1
	for i, t := range view.Tasks {
		if t.Metadata.ID == opt.TaskID {
			taskIndex = i
		}
	}
	if taskIndex == -1 {
		return view, fmt.Errorf("Task not found")
	}
	view.Tasks[taskIndex].Project = newProject
	return view, nil
}
func (h *ViewAPI) handleProjectCreated(ctx context.Context, ev event.Event, view backlogView, logger logger.Logger) (backlogView, error) {
	project := project.Project{}
	err := ev.UnmarshalSpecInto(&project)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to task object: %v", err)
	}
	view.Projects = append(view.Projects, backlogTaskProject{
		ID:   project.Metadata.ID,
		Name: project.Metadata.Name,
	})
	return view, nil
}

func findTask(view backlogView, task string) int {
	taskindex := -1
	for i, t := range view.Tasks {
		if t.Metadata.ID == task {
			taskindex = i
			break
		}
	}
	return taskindex
}
func remove(slice []backlogTask, s int) []backlogTask {
	return append(slice[:s], slice[s+1:]...)
}
