package project

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/project"
	projectCommand "github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// ViewAPI for single project view
	ViewAPI struct {
		TaskRepo    *task.Repo
		ProjectRepo *project.Repo
		DAL         *DAL
	}

	projectView struct {
		project.Project
		Tasks []task.Task `json:"tasks"`
	}
)

// Get build project view
func (h *ViewAPI) Get(c *gin.Context) {
	u := tenant.UserFromContext(c.Request.Context())
	id := c.Param("id")
	view, err := h.DAL.load(c.Request.Context(), u.Metadata.ID, id)
	if err != nil {
		handleError(400, err, c)
		return
	}
	c.JSON(200, view)
}

func (h *ViewAPI) EventHandlers() map[string]handler.EventHandler {
	return map[string]handler.EventHandler{
		"task.deleted":       h,
		"task.completed":     h,
		"task.reopened":      h,
		"project.task-added": h,
		"project.created":    h,
	}
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}

func (h *ViewAPI) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	if ev.Metadata.Name == "project.created" {
		view, err := h.handleProjectCreated(ctx, ev, logger)
		if err != nil {
			return err
		}
		return h.DAL.create(ctx, ev.Tenant.ID, view)
	}
	project, err := h.findProjectIDFromEvent(ctx, ev, logger)
	if err != nil {
		return err
	}
	if project == "" {
		logger.Info("Evenet not related to any known project", "event", ev.Metadata.Name, "aggregator-type", ev.Metadata.AggregatorRoot, "aggregator-id", ev.Metadata.AggregatorID)
		return nil
	}
	current, err := h.DAL.load(ctx, ev.Tenant.ID, project)
	if err != nil {
		return err
	}
	updated, err := h.handlerUpdateEvent(ctx, ev, current, logger)
	if err != nil {
		return err
	}
	return h.DAL.update(ctx, ev.Tenant.ID, project, updated)
}
func (h *ViewAPI) Name() string {
	return "project_view"
}

func (h *ViewAPI) handlerUpdateEvent(ctx context.Context, ev event.Event, view projectView, logger logger.Logger) (projectView, error) {
	switch ev.Metadata.Name {
	case "task.deleted":
		{
			return h.handleTaskDeleted(ctx, ev, view, logger)
		}
	case "project.task-added":
		{
			return h.handleTaskAddedToProject(ctx, ev, view, logger)
		}
	case "task.completed":
		{
			return h.handleTaskStatusChanged(ctx, ev, view, logger)
		}
	case "task.reopened":
		{
			return h.handleTaskStatusChanged(ctx, ev, view, logger)
		}
	}
	return view, fmt.Errorf("Event handler not found")
}

func (h *ViewAPI) findProjectIDFromEvent(ctx context.Context, ev event.Event, logger logger.Logger) (string, error) {
	if ev.Metadata.Name == "task.deleted" || ev.Metadata.Name == "task.completed" || ev.Metadata.Name == "task.reopened" {
		projects, err := h.ProjectRepo.List(project.QueryOptions{
			UserID: ev.Tenant.ID,
		})
		if err != nil {
			return "", err
		}
		projectID := ""
		for _, p := range projects {
			for _, t := range p.Tasks {
				if t == ev.Metadata.AggregatorID {
					projectID = p.Metadata.ID
				}
			}
		}
		return projectID, nil
	}

	if ev.Metadata.Name == "project.task-added" || ev.Metadata.Name == "project.created" {
		return ev.Metadata.AggregatorID, nil
	}
	return "", nil
}

func (h *ViewAPI) handleTaskDeleted(ctx context.Context, ev event.Event, view projectView, logger logger.Logger) (projectView, error) {
	tsk := task.Task{}
	err := ev.UnmarshalSpecInto(&tsk)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to task object: %v", err)
	}
	taskIndex := findTaskIndex(view, tsk.Metadata.ID)
	if taskIndex == -1 {
		return view, fmt.Errorf("Task not found")
	}
	view.Tasks = append(view.Tasks[:taskIndex], view.Tasks[taskIndex+1:]...)
	tasks := []string{}
	for _, t := range view.Tasks {
		tasks = append(tasks, t.Metadata.ID)
	}
	view.Project.Tasks = tasks
	return view, nil
}
func (h *ViewAPI) handleTaskStatusChanged(ctx context.Context, ev event.Event, view projectView, logger logger.Logger) (projectView, error) {
	task, err := h.TaskRepo.Get(ev.Tenant.ID, ev.Metadata.AggregatorID)
	if err != nil {
		return view, err
	}
	taskIndex := findTaskIndex(view, task.Metadata.ID)
	if taskIndex == -1 {
		return view, fmt.Errorf("Task not found")
	}
	view.Tasks[taskIndex] = task
	return view, nil
}
func (h *ViewAPI) handleTaskAddedToProject(ctx context.Context, ev event.Event, view projectView, logger logger.Logger) (projectView, error) {
	opt := projectCommand.AddTasksCommandOptions{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to AddTasksCommandOptions object: %v", err)
	}
	task, err := h.TaskRepo.Get(ev.Tenant.ID, opt.TaskID)
	if err != nil {
		return view, err
	}
	index := findTaskIndex(view, task.Metadata.ID)
	if index != -1 {
		logger.Info("Task already belongs to this project")
		return view, nil
	}
	view.Tasks = append(view.Tasks, task)
	view.Project.Tasks = append(view.Project.Tasks, task.Metadata.ID)
	return view, nil
}
func (h *ViewAPI) handleProjectCreated(ctx context.Context, ev event.Event, logger logger.Logger) (projectView, error) {
	project := project.Project{}
	err := ev.UnmarshalSpecInto(&project)
	if err != nil {
		return projectView{}, fmt.Errorf("Failed to convert event.spec to Project object: %v", err)
	}
	view := projectView{
		Project: project,
		Tasks:   make([]task.Task, 0),
	}
	return view, nil
}

func findTaskIndex(view projectView, task string) int {
	taskIndex := -1
	for i, t := range view.Tasks {
		if t.Metadata.ID == task {
			taskIndex = i
			break
		}
	}
	return taskIndex
}
