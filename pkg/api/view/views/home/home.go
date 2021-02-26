package home

import (
	"context"
	"fmt"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/list"
	listEvents "github.com/peteqproj/peteq/domain/list/event/handler"
	listEventTypes "github.com/peteqproj/peteq/domain/list/event/types"
	"github.com/peteqproj/peteq/domain/project"
	projectEvent "github.com/peteqproj/peteq/domain/project/event/handler"
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

	homeView struct {
		Lists []homeList `json:"lists"`
	}

	homeList struct {
		list.List
		Tasks []homeTask `json:"tasks"`
	}

	homeTask struct {
		task.Task
		Project *repo.Resource `json:"project,omitempty"`
	}
)

// Get builds home view
// @description Home View
// @tags View
// @produce  json
// @success 200 {object} homeView
// @router /q/home [get]
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
		userEventTypes.UserRegistredEvent:         h,
		listEventTypes.ListCreatedEvent:           h,
		listEventTypes.TaskMovedIntoListEvent:     h,
		taskEventTypes.TaskUpdatedEvent:           h,
		taskEventTypes.TaskDeletedEvent:           h,
		projectEventTypes.TaskAddedToProjectEvent: h,
	}
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
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
	return "home_view"
}

func (h *ViewAPI) handlerUpdateEvent(ctx context.Context, ev event.Event, view homeView, logger logger.Logger) (homeView, error) {
	switch ev.Metadata.Name {
	case listEventTypes.ListCreatedEvent:
		{
			return h.handlerListCreated(ctx, ev, view, logger)
		}
	case listEventTypes.TaskMovedIntoListEvent:
		{
			return h.handlerTaskAddedToList(ctx, ev, view, logger)
		}
	case taskEventTypes.TaskUpdatedEvent:
		{
			return h.handlerTaskUpdated(ctx, ev, view, logger)
		}
	case taskEventTypes.TaskDeletedEvent:
		{
			return h.handlerTaskDeleted(ctx, ev, view, logger)
		}
	case projectEventTypes.TaskAddedToProjectEvent:
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
	spec := listEvents.CreatedSpec{}
	if err := ev.UnmarshalSpecInto(&spec); err != nil {
		return view, err
	}
	found := false
	for _, l := range view.Lists {
		if l.Metadata.ID == spec.ID {
			found = true
		}
	}
	if found {
		logger.Info("List already added to view", "list", spec.ID)
		return view, nil
	}

	view.Lists = append(view.Lists, homeList{
		List: list.List{
			Tasks: []string{},
			Metadata: list.Metadata{
				ID:    spec.ID,
				Name:  spec.Name,
				Index: spec.Index,
			},
			Tenant: ev.Tenant,
		},
		Tasks: []homeTask{},
	})
	return view, nil
}
func (h *ViewAPI) handlerTaskAddedToList(ctx context.Context, ev event.Event, view homeView, logger logger.Logger) (homeView, error) {
	spec := listEvents.TaskMovedSpec{}
	if err := ev.UnmarshalSpecInto(&spec); err != nil {
		return view, err
	}
	task, err := h.TaskRepo.Get(ev.Tenant.ID, spec.TaskID)
	if err != nil {
		return view, err
	}
	sourceIndex := -1
	destinationIndex := -1
	for i, l := range view.Lists {
		if spec.Source != "" && l.Metadata.ID == spec.Source {
			sourceIndex = i
			continue
		}

		if spec.Destination != "" && l.Metadata.ID == spec.Destination {
			destinationIndex = i
			continue
		}
	}

	// search if there is reference for task in any project
	projects, err := h.ProjectRepo.List(ctx, repo.ListOptions{})
	if err != nil {
		return view, err
	}

	projectIndex := -1
	taskInProjectIndex := -1
	for i, p := range projects {
		if taskInProjectIndex != -1 {
			break
		}
		if pspec, ok := p.Spec.(project.Spec); ok {
			for j, t := range pspec.Tasks {
				if t == spec.TaskID {
					projectIndex = i
					taskInProjectIndex = j
					break
				}
			}
		}

	}
	var taskProject *repo.Resource
	if taskInProjectIndex != -1 {
		taskProject = projects[projectIndex]
	}

	// If source found, remove task from source
	if sourceIndex != -1 {
		for i, tid := range view.Lists[sourceIndex].Tasks {
			if tid.Task.Metadata.ID == spec.TaskID {
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
	sort.Slice(view.Lists, func(i, j int) bool {
		return view.Lists[i].Metadata.Index < view.Lists[j].Metadata.Index
	})
	return view, nil
}
func (h *ViewAPI) handlerTaskUpdated(ctx context.Context, ev event.Event, view homeView, logger logger.Logger) (homeView, error) {
	spec := taskEvents.UpdatedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}
	listIndex, taskIndex := findTaskInView(view, spec.ID)
	if taskIndex == -1 {
		// task not in lists, no action to do
		return view, nil
	}
	view.Lists[listIndex].Tasks[taskIndex].Task.Metadata.Description = spec.Description
	view.Lists[listIndex].Tasks[taskIndex].Task.Metadata.Name = spec.Name
	return view, nil
}
func (h *ViewAPI) handlerTaskDeleted(ctx context.Context, ev event.Event, view homeView, logger logger.Logger) (homeView, error) {
	spec := taskEvents.DeletedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}

	listIndex, taskIndex := findTaskInView(view, spec.ID)
	if taskIndex == -1 {
		// task not in lists
		return view, nil
	}
	view.Lists[listIndex].Tasks = append(view.Lists[listIndex].Tasks[:taskIndex], view.Lists[listIndex].Tasks[taskIndex+1:]...)
	return view, nil
}
func (h *ViewAPI) handlerTaskAddedToProject(ctx context.Context, ev event.Event, view homeView, logger logger.Logger) (homeView, error) {
	spec := projectEvent.TaskAddedSpec{}
	err := ev.UnmarshalSpecInto(&spec)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to AddTasksCommandOptions object: %v", err)
	}
	newProject, err := h.ProjectRepo.Get(ctx, repo.GetOptions{
		ID: spec.Project,
	})
	if err != nil {
		return view, err
	}
	listIndex, taskIndex := findTaskInView(view, spec.TaskID)
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
			if t.Task.Metadata.ID == id {
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
