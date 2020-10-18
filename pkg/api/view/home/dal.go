package home

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/imdario/mergo"
	"github.com/peteqproj/peteq/domain/list"
	listCommand "github.com/peteqproj/peteq/domain/list/command"
	"github.com/peteqproj/peteq/domain/project"
	projectCommand "github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
)

const dbTableName = "view_home"

type (
	DAL struct {
		DB *sql.DB
	}

	listCreatedHandler struct {
		dal *DAL
	}
	listTaskMovedHandler struct {
		dal         *DAL
		taskRepo    *task.Repo
		projectRepo *project.Repo
	}
	taskCreatedHandler struct {
		dal *DAL
	}
	taskUpdateHandler struct {
		dal *DAL
	}
	taskDeletedHandler struct {
		dal      *DAL
		taskRepo *task.Repo
	}
	userRegistredHandler struct {
		dal *DAL
	}
	projectTaskAddedHandler struct {
		dal         *DAL
		projectRepo *project.Repo
		taskRepo    *task.Repo
	}
)

func (d *DAL) Get(ctx context.Context, user string) (homeView, error) {
	return d.loadHomeView(ctx, user)
}

func (d *DAL) loadHomeView(ctx context.Context, user string) (homeView, error) {
	q, _, err := goqu.From(dbTableName).Where(exp.Ex{
		"userid": []string{user},
	}).ToSQL()
	if err != nil {
		return homeView{}, fmt.Errorf("Failed to build SQL query: %w", err)
	}
	row := d.DB.QueryRowContext(ctx, q)
	view := ""
	userid := ""
	if err := row.Scan(&userid, &view); err != nil {
		return homeView{}, fmt.Errorf("Failed to scan into homeView object: %v", err)
	}
	v := homeView{}
	if err := json.Unmarshal([]byte(view), &v); err != nil {
		return v, err
	}
	return v, nil
}
func (d *DAL) updateView(ctx context.Context, user string, v homeView) error {
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
	rows, err := d.DB.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("Failed to update view_home table: %v", err)
	}
	return rows.Close()
}
func (d *DAL) updateTask(ctx context.Context, user string, task homeTask) error {
	view, err := d.loadHomeView(ctx, user)
	if err != nil {
		return err
	}
	listIndex, taskIndex := findTaskInView(view, task.Metadata.ID)
	if taskIndex == -1 {
		// task not in lists, no action to do
		return nil
	}
	if err := mergo.Merge(&view.Lists[listIndex].Tasks[taskIndex], task, mergo.WithOverwriteWithEmptyValue); err != nil {
		return fmt.Errorf("Failed to update task: %w", err)
	}
	return d.updateView(ctx, user, view)
}
func (l listCreatedHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	curr, err := l.dal.loadHomeView(ctx, e.Tenant.ID)
	if err != nil {
		return err
	}
	opt := listCommand.CreateCommandOptions{}
	if err := e.UnmarshalSpecInto(&opt); err != nil {
		return err
	}
	curr.Lists = append(curr.Lists, homeList{
		List: list.List{
			Tasks: []string{},
			Metadata: list.Metadata{
				ID:    opt.ID,
				Name:  opt.Name,
				Index: opt.Index,
			},
			Tenant: e.Tenant,
		},
		Tasks: []homeTask{},
	})
	return l.dal.updateView(ctx, e.Tenant.ID, curr)
}
func (l listTaskMovedHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	opt := listCommand.MoveTaskArguments{}
	if err := e.UnmarshalSpecInto(&opt); err != nil {
		return err
	}
	view, err := l.dal.loadHomeView(ctx, e.Tenant.ID)
	if err != nil {
		return err
	}
	task, err := l.taskRepo.Get(e.Tenant.ID, opt.TaskID)
	if err != nil {
		return err
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
	projects, err := l.projectRepo.List(project.QueryOptions{
		UserID: e.Tenant.ID,
	})
	if err != nil {
		return err
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
	return l.dal.updateView(ctx, e.Tenant.ID, view)
}
func (t taskUpdateHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	task := task.Task{}
	err := e.UnmarshalSpecInto(&task)
	if err != nil {
		return fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}
	return t.dal.updateTask(ctx, e.Tenant.ID, homeTask{
		Task: task,
	})
}
func (t taskDeletedHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	view, err := t.dal.loadHomeView(ctx, e.Tenant.ID)
	if err != nil {
		return err
	}
	task := task.Task{}
	err = e.UnmarshalSpecInto(&task)
	if err != nil {
		return fmt.Errorf("Failed to convert event.spec to Task object: %v", err)
	}

	listIndex, taskIndex := findTaskInView(view, task.Metadata.ID)
	if taskIndex == -1 {
		// task not in lists
		return nil
	}
	view.Lists[listIndex].Tasks = append(view.Lists[listIndex].Tasks[:taskIndex], view.Lists[listIndex].Tasks[taskIndex+1:]...)
	return t.dal.updateView(ctx, e.Tenant.ID, view)
}
func (t userRegistredHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {

	v := homeView{
		Lists: []homeList{},
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	q, _, err := goqu.Insert(dbTableName).Cols("userid", "info").Vals(goqu.Vals{e.Tenant.ID, string(b)}).ToSQL()
	if err != nil {
		return err
	}
	res, err := t.dal.DB.Query(q)
	if err != nil {
		return err
	}
	defer res.Close()
	return nil
}
func (p projectTaskAddedHandler) Handle(ctx context.Context, e event.Event, logger logger.Logger) error {
	curr, err := p.dal.loadHomeView(ctx, e.Tenant.ID)
	if err != nil {
		return err
	}

	opt := projectCommand.AddTasksCommandOptions{}
	err = e.UnmarshalSpecInto(&opt)
	if err != nil {
		return fmt.Errorf("Failed to convert event.spec to AddTasksCommandOptions object: %v", err)
	}
	newProject, err := p.projectRepo.Get(e.Tenant.ID, opt.Project)
	if err != nil {
		return err
	}
	listIndex, taskIndex := findTaskInView(curr, opt.TaskID)
	if taskIndex == -1 {
		// task not found in lists, not an error
		return nil
	}
	curr.Lists[listIndex].Tasks[taskIndex].Project = newProject
	return p.dal.updateView(ctx, e.Tenant.ID, curr)
}

func (l listCreatedHandler) Name() string {
	return "view_home_listCreatedHandler"
}
func (l listTaskMovedHandler) Name() string {
	return "view_home_listTaskMovedHandler"
}
func (t taskUpdateHandler) Name() string {
	return "view_home_taskUpdateHandler"
}
func (t taskDeletedHandler) Name() string {
	return "view_home_taskDeletedHandler"
}
func (t userRegistredHandler) Name() string {
	return "view_home_userRegistredHandler"
}
func (p projectTaskAddedHandler) Name() string {
	return "view_home_projectTaskAddedHandler"
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
