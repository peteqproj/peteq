package project

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
	"github.com/peteqproj/peteq/pkg/logger"
)

const dbTableName = "view_project"

type (
	DAL struct {
		DB *sql.DB
	}
	taskDeletedHandler struct {
		dal         *DAL
		taskRepo    *task.Repo
		projectRepo *project.Repo
	}
	projectCreatedHandler struct {
		dal *DAL
	}
	projectTaskAddedHandler struct {
		dal         *DAL
		projectRepo *project.Repo
		taskRepo    *task.Repo
	}
	taskStatusChangedHandler struct {
		dal         *DAL
		taskRepo    *task.Repo
		projectRepo *project.Repo
	}
	taskReopenedHandler struct {
		dal         *DAL
		taskRepo    *task.Repo
		projectRepo *project.Repo
	}
)

func (d *DAL) Get(ctx context.Context, user string, project string) (projectView, error) {
	return d.loadprojectView(ctx, user, project)
}

func (d *DAL) loadprojectView(ctx context.Context, user string, project string) (projectView, error) {
	q, _, err := goqu.From(dbTableName).Where(exp.Ex{
		"userid":    user,
		"projectid": project,
	}).ToSQL()
	if err != nil {
		return projectView{}, fmt.Errorf("Failed to build SQL query: %w", err)
	}
	row := d.DB.QueryRowContext(ctx, q)
	view := ""
	userid := ""
	projectid := ""
	if err := row.Scan(&userid, &projectid, &view); err != nil {
		return projectView{}, fmt.Errorf("Failed to scan into projectView object: %v", err)
	}
	v := projectView{}
	if err := json.Unmarshal([]byte(view), &v); err != nil {
		return v, err
	}
	return v, nil
}
func (d *DAL) updateView(ctx context.Context, user string, project string, v projectView) error {
	res, err := json.Marshal(v)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Update(dbTableName).
		Set(goqu.Record{"info": string(res)}).
		Where(exp.Ex{
			"userid":    user,
			"projectid": project,
		}).
		ToSQL()
	if err != nil {
		return err
	}
	rows, err := d.DB.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("Failed to update view_project table: %v", err)
	}
	return rows.Close()
}
func (d *DAL) updateTask(ctx context.Context, user string, project string, task task.Task) error {
	view, err := d.loadprojectView(ctx, user, project)
	if err != nil {
		return err
	}
	taskIndex := findTaskIndex(view, task.Metadata.ID)
	if taskIndex == -1 {
		return fmt.Errorf("Task not found")
	}
	view.Tasks[taskIndex] = task
	return d.updateView(ctx, user, project, view)
}

func (t taskDeletedHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	projects, err := t.projectRepo.List(project.QueryOptions{
		UserID: e.Tenant.ID,
	})
	projectID := ""
	for _, p := range projects {
		for _, t := range p.Tasks {
			if t == e.Metadata.AggregatorID {
				projectID = p.Metadata.ID
			}
		}
	}
	if projectID == "" {
		fmt.Println("Task is not assigned to any project")
		return nil
	}
	view, err := t.dal.loadprojectView(ctx, e.Tenant.ID, projectID)
	if err != nil {
		return err
	}
	tsk := task.Task{}
	err = e.UnmarshalSpecInto(&tsk)
	if err != nil {
		return fmt.Errorf("Failed to convert event.spec to task object: %v", err)
	}
	taskIndex := findTaskIndex(view, tsk.Metadata.ID)
	if taskIndex == -1 {
		return fmt.Errorf("Task not found")
	}
	fmt.Println("Removing", taskIndex)
	view.Tasks = append(view.Tasks[:taskIndex], view.Tasks[taskIndex+1:]...)
	tasks := []string{}
	for _, t := range view.Tasks {
		tasks = append(tasks, t.Metadata.ID)
	}
	view.Project.Tasks = tasks
	return t.dal.updateView(ctx, e.Tenant.ID, projectID, view)
}
func (p projectTaskAddedHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	curr, err := p.dal.loadprojectView(ctx, e.Tenant.ID, e.Metadata.AggregatorID)
	if err != nil {
		return err
	}

	opt := projectCommand.AddTasksCommandOptions{}
	err = e.UnmarshalSpecInto(&opt)
	if err != nil {
		return fmt.Errorf("Failed to convert event.spec to AddTasksCommandOptions object: %v", err)
	}
	task, err := p.taskRepo.Get(e.Tenant.ID, opt.TaskID)
	if err != nil {
		return err
	}
	index := findTaskIndex(curr, task.Metadata.ID)
	if index != -1 {
		fmt.Println("Task already belongs to this project")
		return nil
	}
	curr.Tasks = append(curr.Tasks, task)
	curr.Project.Tasks = append(curr.Project.Tasks, task.Metadata.ID)
	return p.dal.updateView(ctx, e.Tenant.ID, e.Metadata.AggregatorID, curr)
}
func (p projectCreatedHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	project := project.Project{}
	err := e.UnmarshalSpecInto(&project)
	if err != nil {
		return fmt.Errorf("Failed to convert event.spec to Project object: %v", err)
	}
	view := projectView{
		Project: project,
		Tasks:   make([]task.Task, 0),
	}
	b, err := json.Marshal(view)
	if err != nil {
		return err
	}
	q, _, err := goqu.Insert(dbTableName).Cols("userid", "projectid", "info").Vals(goqu.Vals{e.Tenant.ID, project.Metadata.ID, string(b)}).ToSQL()
	if err != nil {
		return err
	}
	_, err = p.dal.DB.Query(q)
	return err
}
func (t taskStatusChangedHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	task, err := t.taskRepo.Get(e.Tenant.ID, e.Metadata.AggregatorID)
	if err != nil {
		return err
	}
	projects, err := t.projectRepo.List(project.QueryOptions{
		UserID: e.Tenant.ID,
	})
	if err != nil {
		return err
	}
	projectID := ""
	for _, p := range projects {
		if projectID != "" {
			break
		}
		for _, t := range p.Tasks {
			if t == task.Metadata.ID {
				projectID = p.Metadata.ID
				break
			}
		}
	}
	if projectID == "" {
		// task not assign to any project
		return nil
	}

	return t.dal.updateTask(ctx, e.Tenant.ID, projectID, task)
}

func (t taskDeletedHandler) Name() string {
	return "view_project_taskDeletedHandler"
}
func (p projectTaskAddedHandler) Name() string {
	return "view_project_projectTaskAddedHandler"
}
func (p projectCreatedHandler) Name() string {
	return "view_project_projectCreatedHandler"
}
func (t taskStatusChangedHandler) Name() string {
	return "view_project_taskStatusChangedHandler"
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
