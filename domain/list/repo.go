package list

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/imdario/mergo"
	"github.com/peteqproj/peteq/pkg/logger"
)

const dbName = "repo_lists"

var errNotFound = errors.New("List not found")

type (
	// Repo is list repository
	// it works on the view db to read/write from it
	Repo struct {
		DB     *sql.DB
		Logger logger.Logger
	}

	// QueryOptions to get task list
	QueryOptions struct {
		UserID string
		noUser bool
	}
)

// List returns set of list
func (r *Repo) List(options QueryOptions) ([]List, error) {
	return r.find(context.Background(), options.UserID)
}

// Get returns list given list id
func (r *Repo) Get(userID string, id string) (List, error) {
	lists, err := r.find(context.Background(), userID, id)
	if err != nil {
		return List{}, err
	}
	if len(lists) == 0 {
		return List{}, errNotFound
	}
	return lists[0], nil
}

// Create will save new task into db
func (r *Repo) Create(user string, l List) error {
	return r.create(context.Background(), user, l)
}

// Delete will remove task from db
func (r *Repo) Delete(userID string, id string) error {
	return r.delete(context.Background(), userID, id)
}

// Update will update given task
func (r *Repo) Update(user string, l List) error {
	list, err := r.Get(user, l.Metadata.ID)
	if err != nil {
		return fmt.Errorf("Failed to read previous task: %w", err)
	}
	if err := mergo.Merge(&list, l, mergo.WithOverwriteWithEmptyValue); err != nil {
		return err
	}
	return r.update(context.Background(), user, list)
}

// MoveTask will move tasks from one list to another one
// TODO: Validation source and destination are exists
func (r *Repo) MoveTask(userID string, sourceID string, destinationID string, task string) error {
	var source *List
	var destination *List
	if sourceID != "" {
		s, err := r.Get(userID, sourceID)
		if err != nil {
			return err
		}
		source = &s
	}
	if destinationID != "" {
		d, err := r.Get(userID, destinationID)
		if err != nil {
			return err
		}
		destination = &d
	}

	// If source found, remove task from source
	if source != nil {
		for i, tid := range source.Tasks {
			if tid == task {
				source.Tasks = remove(source.Tasks, i)
				break
			}
		}
		if err := r.update(context.Background(), userID, *source); err != nil {
			return err
		}
	}

	// If destination found add it to destination
	if destination != nil {
		destination.Tasks = append(destination.Tasks, task)
		if err := r.update(context.Background(), userID, *destination); err != nil {
			return err
		}
	}
	return nil
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func (r *Repo) create(ctx context.Context, user string, list List) error {
	l, err := json.Marshal(list)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Insert(dbName).
		Cols("userid", "listid", "info").
		Vals(goqu.Vals{user, list.Metadata.ID, string(l)}).
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
func (r *Repo) find(ctx context.Context, user string, ids ...string) ([]List, error) {
	e := exp.Ex{
		"userid": user,
	}
	if len(ids) > 0 {
		e = exp.Ex{
			"userid": user,
			"listid": ids,
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
	set := []List{}
	for rows.Next() {
		uid := ""
		lid := ""
		l := ""
		if err := rows.Scan(&uid, &lid, &l); err != nil {
			return nil, err
		}
		prj := List{}
		json.Unmarshal([]byte(l), &prj)
		set = append(set, prj)
	}
	return set, nil
}
func (r *Repo) delete(ctx context.Context, user string, ids ...string) error {
	q, _, err := goqu.
		Delete(dbName).
		Where(exp.Ex{
			"userid": user,
			"listid": ids,
		}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, q)
	return err
}
func (r *Repo) update(ctx context.Context, user string, list List) error {
	l, err := json.Marshal(list)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Update(dbName).
		Where(exp.Ex{
			"userid": user,
			"listid": list.Metadata.ID,
		}).
		Set(goqu.Record{"info": string(l)}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, q)
	return err
}
