package projects

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
	// ViewAPI for projects view
	ViewAPI struct {
		TaskRepo    *task.Repo
		ProjectRepo *project.Repo
		DAL         *DAL
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
func (b *ViewAPI) Get(c *gin.Context) {
	u := tenant.UserFromContext(c.Request.Context())
	view, err := b.DAL.load(c.Request.Context(), u.Metadata.ID)
	if err != nil {
		handleError(400, err, c)
		return
	}
	c.JSON(200, view)
}

func (h *ViewAPI) EventHandlers() map[string]handler.EventHandler {
	return map[string]handler.EventHandler{
		"task.deleted":       h,
		"project.task-added": h,
		"project.created":    h,
		"user.registred":     h,
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
	return "projects_view"
}

func (h *ViewAPI) handlerUpdateEvent(ctx context.Context, ev event.Event, view projectsView, logger logger.Logger) (projectsView, error) {
	switch ev.Metadata.Name {
	case "task.deleted":
		{
			return h.handlerTaskDeleted(ctx, ev, view, logger)
		}
	case "project.task-added":
		{
			return h.handlerTaskAddedToProject(ctx, ev, view, logger)
		}
	case "project.created":
		{
			return h.handlerProjectCreated(ctx, ev, view, logger)
		}
	}
	return view, fmt.Errorf("Event handler not found")
}

func (h *ViewAPI) handlerUserRegistration(ctx context.Context, ev event.Event, logger logger.Logger) error {
	v := projectsView{
		Projects: make([]populatedProject, 0),
	}
	return h.DAL.create(ctx, ev.Tenant.ID, v)
}

func (h *ViewAPI) handlerTaskDeleted(ctx context.Context, ev event.Event, view projectsView, logger logger.Logger) (projectsView, error) {
	task := task.Task{}
	err := ev.UnmarshalSpecInto(&task)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}
	projectIndex, taskIndex := findTaskInView(view, task.Metadata.ID)
	if taskIndex == -1 {
		return view, fmt.Errorf("Task not found")
	}
	view.Projects[projectIndex].Tasks = append(view.Projects[projectIndex].Tasks[:taskIndex], view.Projects[projectIndex].Tasks[taskIndex+1:]...)
	return view, nil
}
func (h *ViewAPI) handlerTaskAddedToProject(ctx context.Context, ev event.Event, view projectsView, logger logger.Logger) (projectsView, error) {
	opt := projectCommand.AddTasksCommandOptions{}
	err := ev.UnmarshalSpecInto(&opt)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to AddTasksCommandOptions object: %v", err)
	}
	newTask, err := h.TaskRepo.Get(ev.Tenant.ID, opt.TaskID)
	if err != nil {
		return view, err
	}
	projectIndex := -1
	for i, p := range view.Projects {
		if p.Metadata.ID == opt.Project {
			projectIndex = i
			break
		}
	}
	if projectIndex == -1 {
		return view, fmt.Errorf("Project not found")
	}
	view.Projects[projectIndex].Tasks = append(view.Projects[projectIndex].Tasks, newTask)
	view.Projects[projectIndex].Project.Tasks = append(view.Projects[projectIndex].Project.Tasks, newTask.Metadata.ID)
	return view, nil
}
func (h *ViewAPI) handlerProjectCreated(ctx context.Context, ev event.Event, view projectsView, logger logger.Logger) (projectsView, error) {
	project := project.Project{}
	err := ev.UnmarshalSpecInto(&project)
	if err != nil {
		return view, fmt.Errorf("Failed to convert event.spec to Project object: %v", err)
	}
	view.Projects = append(view.Projects, populatedProject{Project: project, Tasks: make([]task.Task, 0)})
	return view, nil
}

func findTaskInView(view projectsView, id string) (int, int) {
	projectIndex := -1
	taskIndex := -1
	for i, p := range view.Projects {
		if taskIndex != -1 {
			break
		}
		for j, t := range p.Tasks {
			if t.Metadata.ID == id {
				projectIndex = i
				taskIndex = j
				break
			}
		}
	}
	return projectIndex, taskIndex
}
