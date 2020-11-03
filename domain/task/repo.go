package task

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/imdario/mergo"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/logger"
)

const dbName = "repo_tasks"

var errNotFound = errors.New("Task not found")

type (
	// Repo is task repository
	// it works on the view db to read/write from it
	Repo struct {
		DB     db.Database
		Logger logger.Logger
	}

	// ListOptions to get task list
	ListOptions struct {
		UserID string
	}
)

// List returns list of tasks
func (r *Repo) List(options ListOptions) ([]Task, error) {
	return r.find(context.Background(), options.UserID)
}

// Get returns task given task id
func (r *Repo) Get(userID string, id string) (Task, error) {
	tasks, err := r.find(context.Background(), userID, id)
	if err != nil {
		return Task{}, err
	}
	if len(tasks) == 0 {
		return Task{}, errNotFound
	}
	return tasks[0], nil
}

// Create will save new task into db
func (r *Repo) Create(user string, t Task) error {
	return r.create(context.Background(), user, t)
}

// Delete will remove task from db
func (r *Repo) Delete(userID string, id string) error {
	return r.delete(context.Background(), userID, id)
}

// Update will update given task
func (r *Repo) Update(user string, newTask Task) error {
	t, err := r.Get(user, newTask.Metadata.ID)
	if err != nil {
		return err
	}
	if err := mergo.Merge(&t, newTask, mergo.WithOverwriteWithEmptyValue); err != nil {
		return err
	}
	return r.update(context.Background(), user, t)
}

func (r *Repo) create(ctx context.Context, user string, task Task) error {
	t, err := json.Marshal(task)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Insert(dbName).
		Cols("userid", "taskid", "info").
		Vals(goqu.Vals{user, task.Metadata.ID, string(t)}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, q)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) find(ctx context.Context, user string, ids ...string) ([]Task, error) {
	e := exp.Ex{
		"userid": user,
	}
	if len(ids) > 0 {
		e = exp.Ex{
			"userid": user,
			"taskid": ids,
		}
	}
	q, _, err := goqu.
		From(dbName).
		Where(e).
		ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	set := []Task{}
	for rows.Next() {
		uid := ""
		tid := ""
		t := ""
		if err := rows.Scan(&uid, &tid, &t); err != nil {
			return nil, err
		}
		tsk := Task{}
		json.Unmarshal([]byte(t), &tsk)
		set = append(set, tsk)
	}
	return set, rows.Close()
}
func (r *Repo) delete(ctx context.Context, user string, ids ...string) error {
	q, _, err := goqu.
		Delete(dbName).
		Where(exp.Ex{
			"userid": user,
			"taskid": ids,
		}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, q)
	return err
}
func (r *Repo) update(ctx context.Context, user string, task Task) error {
	t, err := json.Marshal(task)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Update(dbName).
		Where(exp.Ex{
			"userid": user,
			"taskid": task.Metadata.ID,
		}).
		Set(goqu.Record{"info": string(t)}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, q)
	return err
}
