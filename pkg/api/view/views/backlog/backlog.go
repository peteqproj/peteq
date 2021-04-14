package backlog

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/list"
	listEvents "github.com/peteqproj/peteq/domain/list/event/handler"
	listEventTypes "github.com/peteqproj/peteq/domain/list/event/types"
	projectEvents "github.com/peteqproj/peteq/domain/project/event/handler"
	projectEventTypes "github.com/peteqproj/peteq/domain/project/event/types"
	"github.com/peteqproj/peteq/domain/task"
	taskEvents "github.com/peteqproj/peteq/domain/task/event/handler"
	taskEventTypes "github.com/peteqproj/peteq/domain/task/event/types"
	userEventTypes "github.com/peteqproj/peteq/domain/user/event/types"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/event/handler"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/repo"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// ViewAPI for backlog view
	ViewAPI struct {
		TaskRepo    *task.Repo
		ListRepo    *list.Repo
		ProjectRepo *repo.Repo
		DAL         *DAL
	}

	backlogView struct {
		Tasks    []backlogTask        `json:"tasks"`
		Lists    []backlogTaskList    `json:"lists"`
		Projects []backlogTaskProject `json:"projects"`
	}

	backlogTask struct {
		Task    task.Task          `json:"task"`
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
// @description Backlog View
// @tags View
// @produce  json
// @success 200 {object} backlogView
// @router /q/backlog [get]
// @Security ApiKeyAuth
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
		listEventTypes.TaskMovedIntoListEvent:     h,
		taskEventTypes.TaskCreatedEvent:           h,
		taskEventTypes.TaskUpdatedEvent:           h,
		taskEventTypes.TaskStatusChanged:          h,
		taskEventTypes.TaskDeletedEvent:           h,
		userEventTypes.UserRegistredEvent:         h,
		projectEventTypes.ProjectCreatedEvent:     h,
		projectEventTypes.TaskAddedToProjectEvent: h,
	}
}

func (h *ViewAPI) Handle(ctx context.Context, ev event.Event, logger logger.Logger) error {
	if ev.Metadata.Name == userEventTypes.UserRegistredEvent {
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
	case listEventTypes.TaskMovedIntoListEvent:
		{
			return h.handleTaskMovedToList(ctx, ev, view, logger)
		}
	case taskEventTypes.TaskCreatedEvent:
		{
			return h.handleTaskCreated(ctx, ev, view, logger)
		}
	case taskEventTypes.TaskUpdatedEvent:
		{
			return h.handleTaskUpdated(ctx, ev, view, logger)
		}
	case taskEventTypes.TaskStatusChanged:
		{
			return h.handleTaskStatusChanged(ctx, ev, view, logger)
		}
	case taskEventTypes.TaskDeletedEvent:
		{
			return h.handleTaskDeleted(ctx, ev, view, logger)
		}
	case projectEventTypes.ProjectCreatedEvent:
		{
			return h.handleProjectCreated(ctx, ev, view, logger)
		}
	case projectEventTypes.TaskAddedToProjectEvent:
		{
			return h.handleTaskAddedToProject(ctx, ev, view, logger)
		}
	}
	return view, fmt.Errorf("Event handler not found")
}

func (h *ViewAPI) handleTaskMovedToList(ctx context.Context, ev event.Event, view backlogView, logger logger.Logger) (backlogView, error) {
	spec := listEvents.TaskMovedSpec{}
	if err := ev.UnmarshalSpecInto(&spec); err != nil {
		return view, err
	}
	var newList *list.List
	if spec.Destination != "" {
		list, err := h.ListRepo.GetById(ctx, spec.Destination)
		if err != nil {
			return view, err
		}
		logger.Info("Destination is set", "id", spec.Destination, "name", list.Metadata.Name)
		newList = list
	}
	taskIndex := -1
	for i, t := range view.Tasks {
		if t.Task.Metadata.ID == spec.TaskID {
			logger.Info("Task found in view", "id", spec.TaskID, "index", i)
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
	spec := taskEvents.CreatedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}
	t := task.Task{
		Metadata: task.Metadata{
			ID:   spec.ID,
			Name: spec.Name,
			// Description: &spec.Description,
		},
	}
	view.Tasks = append(view.Tasks, backlogTask{
		Task: t,
	})
	return view, nil
}
func (h *ViewAPI) handleTaskUpdated(ctx context.Context, ev event.Event, view backlogView, logger logger.Logger) (backlogView, error) {
	spec := taskEvents.UpdatedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}
	index := findTask(view, spec.ID)
	if index == -1 {
		return view, fmt.Errorf("Task not found")
	}
	view.Tasks[index].Task.Metadata.Name = spec.Name
	// view.Tasks[index].Task.Metadata.Description = &spec.Description
	return view, nil
}
func (h *ViewAPI) handleTaskStatusChanged(ctx context.Context, ev event.Event, view backlogView, logger logger.Logger) (backlogView, error) {
	spec := taskEvents.StatusChangedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to StatusChangedSpec object: %v", err)
	}
	index := findTask(view, ev.Metadata.AggregatorID)
	if index == -1 {
		return view, fmt.Errorf("Task not found")
	}
	view.Tasks[index].Task.Spec.Completed = true
	return view, nil
}
func (h *ViewAPI) handleTaskDeleted(ctx context.Context, ev event.Event, view backlogView, logger logger.Logger) (backlogView, error) {
	spec := taskEvents.DeletedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}

	taskIndex := -1
	for i, t := range view.Tasks {
		if t.Task.Metadata.ID == spec.ID {
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
	spec := projectEvents.TaskAddedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to task object: %v", err)
	}
	newProject := backlogTaskProject{}
	if spec.Project != "" {
		prj, err := h.ProjectRepo.Get(ctx, repo.GetOptions{
			ID: spec.Project,
		})
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
		if t.Task.Metadata.ID == spec.TaskID {
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
	spec := projectEvents.CreatedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to task object: %v", err)
	}
	view.Projects = append(view.Projects, backlogTaskProject{
		ID:   spec.ID,
		Name: spec.Name,
	})
	return view, nil
}

func findTask(view backlogView, task string) int {
	taskindex := -1
	for i, t := range view.Tasks {
		if t.Task.Metadata.ID == task {
			taskindex = i
			break
		}
	}
	return taskindex
}
func remove(slice []backlogTask, s int) []backlogTask {
	return append(slice[:s], slice[s+1:]...)
}
