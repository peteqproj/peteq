package backlog

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/peteqproj/peteq/domain/list"
	listCommand "github.com/peteqproj/peteq/domain/list/command"
	"github.com/peteqproj/peteq/domain/project"
	projectCommand "github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

const dbTableName = "view_backlog"

type (
	DAL struct {
		DB *sql.DB
	}

	taskCreatedHandler struct {
		dal *DAL
	}
	taskUpdateHandler struct {
		dal *DAL
	}
	taskStatusChangedHandler struct {
		dal      *DAL
		taskRepo *task.Repo
	}
	taskDeletedHandler struct {
		dal *DAL
	}
	userRegistredHandler struct {
		dal *DAL
	}
	listTaskMovedHandler struct {
		dal      *DAL
		taskRepo *task.Repo
		listRepo *list.Repo
	}
	projectCreatedHandler struct {
		dal         *DAL
		projectRepo *project.Repo
	}
	projectTaskAddedHandler struct {
		dal         *DAL
		projectRepo *project.Repo
		taskRepo    *task.Repo
	}
)

func (d *DAL) Get(ctx context.Context, user string) (backlogView, error) {
	return d.loadBacklog(ctx, user)
}
func (d *DAL) loadBacklog(ctx context.Context, user string) (backlogView, error) {
	q, _, err := goqu.From(dbTableName).Where(exp.Ex{
		"userid": []string{user},
	}).ToSQL()
	if err != nil {
		return backlogView{}, fmt.Errorf("Failed to build SQL query: %w", err)
	}
	row := d.DB.QueryRowContext(ctx, q)

	view := ""
	userid := ""
	if err := row.Scan(&userid, &view); err != nil {
		return backlogView{}, fmt.Errorf("Failed to scan into backlogView object: %v", err)
	}
	v := backlogView{}
	if err := json.Unmarshal([]byte(view), &v); err != nil {
		return v, err
	}
	return v, nil
}
func (d *DAL) updateTask(ctx context.Context, user string, task backlogTask) error {
	curr, err := d.loadBacklog(ctx, user)
	if err != nil {
		return err
	}
	taskindex := findTask(curr, task.Task.Metadata.ID)
	if taskindex == -1 {
		return fmt.Errorf("Task not found")
	}
	curr.Tasks[taskindex] = task
	return d.updateView(ctx, user, curr)
}
func (d *DAL) updateView(ctx context.Context, user string, view backlogView) error {
	res, err := json.Marshal(view)
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
	_, err = d.DB.Query(q)
	if err != nil {
		return fmt.Errorf("Failed to update db: %v", err)
	}
	return nil
}

func (t taskCreatedHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	curr, err := t.dal.loadBacklog(ctx, e.Tenant.ID)
	if err != nil {
		return err
	}

	task := task.Task{}
	err = e.UnmarshalSpecInto(&task)
	if err != nil {
		return fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}
	curr.Tasks = append(curr.Tasks, backlogTask{
		Task: task,
	})
	res, err := json.Marshal(curr)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Update(dbTableName).
		Set(goqu.Record{"info": string(res)}).
		Where(exp.Ex{
			"userid": []string{e.Tenant.ID},
		}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = t.dal.DB.Query(q)
	if err != nil {
		return fmt.Errorf("Failed to update db: %v", err)
	}
	return nil
}
func (t taskUpdateHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	view, err := t.dal.loadBacklog(ctx, e.Tenant.ID)
	if err != nil {
		return err
	}
	task := task.Task{}
	err = e.UnmarshalSpecInto(&task)
	if err != nil {
		return fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}
	index := findTask(view, task.Metadata.ID)
	if index == -1 {
		return fmt.Errorf("Task not found")
	}
	view.Tasks[index].Task = task
	return t.dal.updateView(ctx, e.Tenant.ID, view)
}
func (t taskStatusChangedHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	view, err := t.dal.loadBacklog(ctx, e.Tenant.ID)
	if err != nil {
		return err
	}
	task, err := t.taskRepo.Get(e.Tenant.ID, e.Metadata.AggregatorID)
	if err != nil {
		return err
	}
	index := findTask(view, task.Metadata.ID)
	if index == -1 {
		return fmt.Errorf("Task not found")
	}
	view.Tasks[index].Task = task
	return t.dal.updateView(ctx, e.Tenant.ID, view)
}
func (t taskDeletedHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	view, err := t.dal.loadBacklog(ctx, e.Tenant.ID)
	if err != nil {
		return err
	}
	task := task.Task{}
	err = e.UnmarshalSpecInto(&task)
	if err != nil {
		return fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}

	taskIndex := -1
	for i, t := range view.Tasks {
		if t.Metadata.ID == task.Metadata.ID {
			taskIndex = i
			break
		}
	}
	if taskIndex == -1 {
		return fmt.Errorf("Task not found")
	}
	view.Tasks = remove(view.Tasks, taskIndex)
	return t.dal.updateView(ctx, e.Tenant.ID, view)
}
func (u userRegistredHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	v := backlogView{
		Tasks:    make([]backlogTask, 0),
		Lists:    make([]backlogTaskList, 0),
		Projects: make([]backlogTaskProject, 0),
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
func (l listTaskMovedHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	opt := listCommand.MoveTaskArguments{}
	if err := e.UnmarshalSpecInto(&opt); err != nil {
		return err
	}
	view, err := l.dal.loadBacklog(ctx, e.Tenant.ID)
	if err != nil {
		return err
	}
	newList := list.List{}
	if opt.Destination != "" {
		list, err := l.listRepo.Get(e.Tenant.ID, opt.Destination)
		if err != nil {
			return err
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
		return fmt.Errorf("Task not found")
	}
	view.Tasks[taskIndex].List = backlogTaskList{
		ID:   newList.Metadata.ID,
		Name: newList.Metadata.Name,
	}
	return l.dal.updateView(ctx, e.Tenant.ID, view)
}
func (p projectCreatedHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	curr, err := p.dal.loadBacklog(ctx, e.Tenant.ID)
	if err != nil {
		return err
	}

	project := project.Project{}
	err = e.UnmarshalSpecInto(&project)
	if err != nil {
		return fmt.Errorf("Failed to convert event.spec to task object: %v", err)
	}
	curr.Projects = append(curr.Projects, backlogTaskProject{
		ID:   project.Metadata.ID,
		Name: project.Metadata.Name,
	})
	return p.dal.updateView(ctx, e.Tenant.ID, curr)
}
func (p projectTaskAddedHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	curr, err := p.dal.loadBacklog(ctx, e.Tenant.ID)
	if err != nil {
		return err
	}

	opt := projectCommand.AddTasksCommandOptions{}
	err = e.UnmarshalSpecInto(&opt)
	if err != nil {
		return fmt.Errorf("Failed to convert event.spec to task object: %v", err)
	}
	newProject := backlogTaskProject{}
	if opt.Project != "" {
		prj, err := p.projectRepo.Get(e.Tenant.ID, opt.Project)
		if err != nil {
			return err
		}
		newProject = backlogTaskProject{
			ID:   prj.Metadata.ID,
			Name: prj.Metadata.Name,
		}
	}
	for _, t := range curr.Tasks {
		if t.Metadata.ID == opt.TaskID {
			t.Project = newProject
			return p.dal.updateTask(ctx, e.Tenant.ID, t)
		}
	}
	return nil
}

func (t taskCreatedHandler) Name() string {
	return "view_backlog_taskCreatedHandler"
}
func (t taskUpdateHandler) Name() string {
	return "view_backlog_taskUpdateHandler"
}
func (t taskStatusChangedHandler) Name() string {
	return "view_backlog_taskStatusChangedHandler"
}
func (t taskDeletedHandler) Name() string {
	return "view_backlog_taskDeletedHandler"
}
func (u userRegistredHandler) Name() string {
	return "view_backlog_userRegistredHandler"
}
func (l listTaskMovedHandler) Name() string {
	return "view_backlog_listTaskMovedHandler"
}
func (p projectCreatedHandler) Name() string {
	return "view_backlog_projectCreatedHandler"
}
func (p projectTaskAddedHandler) Name() string {
	return "view_backlog_projectTaskAddedHandler"
}

func remove(slice []backlogTask, s int) []backlogTask {
	return append(slice[:s], slice[s+1:]...)
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
