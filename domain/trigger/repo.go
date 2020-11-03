package trigger

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"

	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/logger"
)

const dbName = "repo_triggers"

var errNotFound = errors.New("Trigger not foud")

type (
	// Repo is trigger repository
	// it works on the view db to read/write from it
	Repo struct {
		DB     db.Database
		Logger logger.Logger
	}

	// QueryOptions to get trigger trigger
	QueryOptions struct {
		UserID string
	}
)

// List returns set of trigger
func (r *Repo) List(options QueryOptions) ([]Trigger, error) {
	return r.find(context.Background(), options.UserID)
}

// Get returns trigger given trigger id
func (r *Repo) Get(userID string, id string) (Trigger, error) {
	triggers, err := r.find(context.Background(), userID, id)
	if err != nil {
		return Trigger{}, err
	}
	if len(triggers) == 0 {
		return Trigger{}, errNotFound
	}
	return triggers[0], nil
}

// Create will save new trigger into db
func (r *Repo) Create(user string, trigger Trigger) error {
	return r.create(context.Background(), user, trigger)
}

func (r *Repo) create(ctx context.Context, user string, trigger Trigger) error {
	t, err := json.Marshal(trigger)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Insert(dbName).
		Cols("userid", "triggerid", "info").
		Vals(goqu.Vals{user, trigger.Metadata.ID, string(t)}).
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
func (r *Repo) find(ctx context.Context, user string, ids ...string) ([]Trigger, error) {
	e := exp.Ex{
		"userid": user,
	}
	if len(ids) > 0 {
		e = exp.Ex{
			"userid":    user,
			"triggerid": ids,
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
	set := []Trigger{}
	for rows.Next() {
		uid := ""
		pid := ""
		t := ""
		if err := rows.Scan(&uid, &pid, &t); err != nil {
			return nil, err
		}
		trigger := Trigger{}
		json.Unmarshal([]byte(t), &trigger)
		set = append(set, trigger)
	}
	return set, rows.Close()
}
