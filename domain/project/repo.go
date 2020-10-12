package project

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/imdario/mergo"
	"github.com/peteqproj/peteq/pkg/logger"
)

const dbName = "repo_projects"

var errNotFound = errors.New("Project not foud")

type (
	// Repo is project repository
	// it works on the view db to read/write from it
	Repo struct {
		DB     *sql.DB
		Logger logger.Logger
	}

	// QueryOptions to get project project
	QueryOptions struct {
		UserID string
		noUser bool
	}
)

// List returns set of project
func (r *Repo) List(options QueryOptions) ([]Project, error) {
	return r.find(context.Background(), options.UserID)
}

// Get returns project given project id
func (r *Repo) Get(userID string, id string) (Project, error) {
	projects, err := r.find(context.Background(), userID, id)
	if err != nil {
		return Project{}, err
	}
	if len(projects) == 0 {
		return Project{}, errNotFound
	}
	return projects[0], nil
}

// Create will save new project into db
func (r *Repo) Create(user string, project Project) error {
	return r.create(context.Background(), user, project)
}

// Delete will remove project from db
func (r *Repo) Delete(userID string, id string) error {
	return r.delete(context.Background(), userID, id)
}

// Update will update given project
func (r *Repo) Update(user string, p Project) error {
	prj, err := r.Get(user, p.Metadata.ID)
	if err != nil {
		return err
	}
	if err := mergo.Merge(&prj, p, mergo.WithOverwriteWithEmptyValue); err != nil {
		return err
	}
	return r.update(context.Background(), user, prj)
}

// AddTask adds task to project
// TODO: check that task is not assigned to other project
func (r *Repo) AddTask(userID string, project string, task string) error {
	proj, err := r.Get(userID, project)
	if err != nil {
		return err
	}
	proj.Tasks = append(proj.Tasks, task)
	return r.Update(userID, proj)
}

func (r *Repo) create(ctx context.Context, user string, prj Project) error {
	t, err := json.Marshal(prj)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Insert(dbName).
		Cols("userid", "projectid", "info").
		Vals(goqu.Vals{user, prj.Metadata.ID, string(t)}).
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
func (r *Repo) find(ctx context.Context, user string, ids ...string) ([]Project, error) {
	e := exp.Ex{
		"userid": user,
	}
	if len(ids) > 0 {
		e = exp.Ex{
			"userid":    user,
			"projectid": ids,
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
	set := []Project{}
	for rows.Next() {
		uid := ""
		pid := ""
		t := ""
		if err := rows.Scan(&uid, &pid, &t); err != nil {
			return nil, err
		}
		prj := Project{}
		json.Unmarshal([]byte(t), &prj)
		set = append(set, prj)
	}
	return set, nil
}
func (r *Repo) delete(ctx context.Context, user string, ids ...string) error {
	q, _, err := goqu.
		Delete(dbName).
		Where(exp.Ex{
			"userid":    user,
			"projectid": ids,
		}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, q)
	return err
}
func (r *Repo) update(ctx context.Context, user string, project Project) error {
	t, err := json.Marshal(project)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Update(dbName).
		Where(exp.Ex{
			"userid":    user,
			"projectid": project.Metadata.ID,
		}).
		Set(goqu.Record{"info": string(t)}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, q)
	return err
}
