package project

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/project"
	projectEvent "github.com/peteqproj/peteq/domain/project/event/handler"
	projectEventTypes "github.com/peteqproj/peteq/domain/project/event/types"
	"github.com/peteqproj/peteq/domain/task"
	taskEvents "github.com/peteqproj/peteq/domain/task/event/handler"
	taskEventTypes "github.com/peteqproj/peteq/domain/task/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/repo"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// ViewAPI for single project view
	ViewAPI struct {
		TaskRepo    *task.Repo
		ProjectRepo *repo.Repo
		DAL         *DAL
	}

	projectView struct {
		Project project.Project `json:"project"`
		Tasks   []task.Task     `json:"tasks"`
	}
)

// Get build project view
// @description Single Project View
// @tags View
// @produce  json
// @success 200 {object} projectView
// @Param id path string true "Project ID"
// @router /q/project/{id} [get]
// @Security ApiKeyAuth
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
		taskEventTypes.TaskDeletedEvent:           h,
		taskEventTypes.TaskStatusChanged:          h,
		taskEventTypes.TaskUpdatedEvent:           h,
		projectEventTypes.ProjectCreatedEvent:     h,
		projectEventTypes.TaskAddedToProjectEvent: h,
	}
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}

func (h *ViewAPI) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	if ev.Metadata.Name == projectEventTypes.ProjectCreatedEvent {
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
		logger.Info("Event is not related to any known project", "event", ev.Metadata.Name, "aggregator-type", ev.Metadata.AggregatorRoot, "aggregator-id", ev.Metadata.AggregatorID)
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
	case taskEventTypes.TaskDeletedEvent:
		{
			return h.handleTaskDeleted(ctx, ev, view, logger)
		}
	case projectEventTypes.TaskAddedToProjectEvent:
		{
			return h.handleTaskAddedToProject(ctx, ev, view, logger)
		}
	case taskEventTypes.TaskStatusChanged:
		{
			return h.handleTaskStatusChanged(ctx, ev, view, logger)
		}
	case taskEventTypes.TaskUpdatedEvent:
		{
			return h.handleTaskUpdated(ctx, ev, view, logger)
		}
	}
	return view, fmt.Errorf("Event handler not found")
}

func (h *ViewAPI) findProjectIDFromEvent(ctx context.Context, ev event.Event, logger logger.Logger) (string, error) {
	if ev.Metadata.Name == taskEventTypes.TaskDeletedEvent || ev.Metadata.Name == taskEventTypes.TaskStatusChanged || ev.Metadata.Name == taskEventTypes.TaskUpdatedEvent {
		projects, err := h.ProjectRepo.List(ctx, repo.ListOptions{})
		if err != nil {
			return "", err
		}
		projectID := ""
		for _, p := range projects {
			if spec, ok := p.Spec.(project.Spec); ok {
				for _, t := range spec.Tasks {
					if t == ev.Metadata.AggregatorID {
						projectID = p.Metadata.ID
					}
				}
			}
		}
		return projectID, nil
	}

	if ev.Metadata.Name == projectEventTypes.TaskAddedToProjectEvent || ev.Metadata.Name == projectEventTypes.ProjectCreatedEvent {
		return ev.Metadata.AggregatorID, nil
	}
	return "", nil
}

func (h *ViewAPI) handleTaskDeleted(ctx context.Context, ev event.Event, view projectView, logger logger.Logger) (projectView, error) {
	spec := taskEvents.DeletedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to task object: %v", err)
	}
	taskIndex := findTaskIndex(view, spec.ID)
	if taskIndex == -1 {
		return view, fmt.Errorf("Task not found")
	}
	view.Tasks = append(view.Tasks[:taskIndex], view.Tasks[taskIndex+1:]...)
	return view, nil
}
func (h *ViewAPI) handleTaskStatusChanged(ctx context.Context, ev event.Event, view projectView, logger logger.Logger) (projectView, error) {
	spec := taskEvents.StatusChangedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to StatusChangedSpec object: %v", err)
	}
	taskIndex := findTaskIndex(view, ev.Metadata.AggregatorID)
	if taskIndex == -1 {
		return view, fmt.Errorf("Task not found")
	}
	view.Tasks[taskIndex].Spec.Completed = spec.Completed
	return view, nil
}
func (h *ViewAPI) handleTaskUpdated(ctx context.Context, ev event.Event, view projectView, logger logger.Logger) (projectView, error) {
	spec := taskEvents.UpdatedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to StatusChangedSpec object: %v", err)
	}
	taskIndex := findTaskIndex(view, ev.Metadata.AggregatorID)
	if taskIndex == -1 {
		return view, fmt.Errorf("Task not found")
	}
	if spec.Description != "" {
		view.Tasks[taskIndex].Metadata.Description = utils.PtrString(spec.Description)
	}
	if spec.Name != "" {
		view.Tasks[taskIndex].Metadata.Name = spec.Name
	}
	return view, nil
}
func (h *ViewAPI) handleTaskAddedToProject(ctx context.Context, ev event.Event, view projectView, logger logger.Logger) (projectView, error) {
	spec := projectEvent.TaskAddedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to AddTasksCommandOptions object: %v", err)
	}
	task, err := h.TaskRepo.GetById(ctx, spec.TaskID)
	if err != nil {
		return view, err
	}
	index := findTaskIndex(view, task.Metadata.ID)
	if index != -1 {
		logger.Info("Task already belongs to this project")
		return view, nil
	}
	view.Tasks = append(view.Tasks, *task)
	return view, nil
}
func (h *ViewAPI) handleProjectCreated(ctx context.Context, ev event.Event, logger logger.Logger) (projectView, error) {
	spec := projectEvent.CreatedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return projectView{}, fmt.Errorf("Failed to convert event.spec to Project object: %v", err)
	}
	p := project.Project{
		Metadata: project.Metadata{
			ID:          spec.ID,
			Name:        spec.Name,
			Description: &spec.Description,
		},
		Spec: project.Spec{
			Color:    &spec.Color,
			ImageURL: &spec.ImageURL,
			Tasks:    []string{},
		},
	}
	view := projectView{
		Project: p,
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
