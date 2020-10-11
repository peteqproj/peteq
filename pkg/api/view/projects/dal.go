package projects

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/peteqproj/peteq/domain/project"
	projectCommand "github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
)

const dbTableName = "view_projects"

type (
	DAL struct {
		DB *sql.DB
	}
	taskDeletedHandler struct {
		dal      *DAL
		taskRepo *task.Repo
	}
	projectCreatedHandler struct {
		dal *DAL
	}
	projectTaskAddedHandler struct {
		dal         *DAL
		projectRepo *project.Repo
		taskRepo    *task.Repo
	}
	userRegistredHandler struct {
		dal *DAL
	}
)

func (d *DAL) Get(ctx context.Context, user string) (projectsView, error) {
	return d.loadProjectsView(ctx, user)
}

func (d *DAL) loadProjectsView(ctx context.Context, user string) (projectsView, error) {
	q, _, err := goqu.From(dbTableName).Where(exp.Ex{
		"userid": []string{user},
	}).ToSQL()
	if err != nil {
		return projectsView{}, fmt.Errorf("Failed to build SQL query: %w", err)
	}
	row := d.DB.QueryRowContext(ctx, q)
	view := ""
	userid := ""
	if err := row.Scan(&userid, &view); err != nil {
		return projectsView{}, fmt.Errorf("Failed to scan into projectsView object: %v", err)
	}
	v := projectsView{}
	if err := json.Unmarshal([]byte(view), &v); err != nil {
		return v, err
	}
	return v, nil
}
func (d *DAL) updateView(ctx context.Context, user string, v projectsView) error {
	res, err := json.Marshal(v)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Update(dbTableName).
		Set(goqu.Record{"info": string(res)}).
		Where(exp.Ex{
			"userid": []string{user},
		}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = d.DB.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("Failed to update view_home table: %v", err)
	}
	return nil
}

func (t taskDeletedHandler) Handle(e event.Event) error {
	view, err := t.dal.loadProjectsView(context.Background(), e.Tenant.ID)
	if err != nil {
		return err
	}
	task := task.Task{}
	err = e.UnmarshalSpecInto(&task)
	if err != nil {
		return fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}
	projectIndex, taskIndex := findTaskInView(view, task.Metadata.ID)
	if taskIndex == -1 {
		return fmt.Errorf("Task not found")
	}
	view.Projects[projectIndex].Tasks = append(view.Projects[projectIndex].Tasks[:taskIndex], view.Projects[projectIndex].Tasks[taskIndex+1:]...)
	return t.dal.updateView(context.Background(), e.Tenant.ID, view)
}
func (p projectTaskAddedHandler) Handle(e event.Event) error {
	curr, err := p.dal.loadProjectsView(context.Background(), e.Tenant.ID)
	if err != nil {
		return err
	}

	opt := projectCommand.AddTasksCommandOptions{}
	err = e.UnmarshalSpecInto(&opt)
	if err != nil {
		return fmt.Errorf("Failed to convert event.spec to AddTasksCommandOptions object: %v", err)
	}
	newTask, err := p.taskRepo.Get(e.Tenant.ID, opt.TaskID)
	if err != nil {
		return err
	}
	projectIndex := -1
	for i, p := range curr.Projects {
		if p.Metadata.ID == opt.Project {
			projectIndex = i
			break
		}
	}
	if projectIndex == -1 {
		return fmt.Errorf("Project not found")
	}
	curr.Projects[projectIndex].Tasks = append(curr.Projects[projectIndex].Tasks, newTask)
	curr.Projects[projectIndex].Project.Tasks = append(curr.Projects[projectIndex].Project.Tasks, newTask.Metadata.ID)
	return p.dal.updateView(context.Background(), e.Tenant.ID, curr)
}
func (p projectCreatedHandler) Handle(e event.Event) error {
	curr, err := p.dal.loadProjectsView(context.Background(), e.Tenant.ID)
	if err != nil {
		return err
	}

	project := project.Project{}
	err = e.UnmarshalSpecInto(&project)
	if err != nil {
		return fmt.Errorf("Failed to convert event.spec to Project object: %v", err)
	}
	curr.Projects = append(curr.Projects, populatedProject{Project: project, Tasks: make([]task.Task, 0)})
	return p.dal.updateView(context.Background(), e.Tenant.ID, curr)
}
func (u userRegistredHandler) Handle(e event.Event) error {
	v := projectsView{
		Projects: []populatedProject{},
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	q, _, err := goqu.Insert(dbTableName).Cols("userid", "info").Vals(goqu.Vals{e.Tenant.ID, string(b)}).ToSQL()
	if err != nil {
		return err
	}
	_, err = u.dal.DB.Query(q)
	return err
}

func (t taskDeletedHandler) Name() string {
	return "view_projects_taskDeletedHandler"
}
func (p projectTaskAddedHandler) Name() string {
	return "view_projects_projectTaskAddedHandler"
}
func (p projectCreatedHandler) Name() string {
	return "view_projects_projectCreatedHandler"
}
func (u userRegistredHandler) Name() string {
	return "view_projects_userRegistredHandler"
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
