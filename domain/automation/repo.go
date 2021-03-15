package automation

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"

	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/logger"
)

const (
	dbName               = "repo_automations"
	triggerBindingDbName = "repo_trigger_bindings"
)

var errNotFound = errors.New("Automation not foud")
var errTBNotFound = errors.New("TriggerBinding not foud")

type (
	// Repo is automation repository
	// it works on the view db to read/write from it
	Repo struct {
		DB     db.Database
		Logger logger.Logger
	}

	// QueryOptions to get automation automation
	QueryOptions struct {
		UserID string
	}
)

// List returns set of automation
func (r *Repo) List(options QueryOptions) ([]Automation, error) {
	return r.find(context.Background(), options.UserID)
}

// Get returns automation given automation id
func (r *Repo) Get(userID string, id string) (Automation, error) {
	automations, err := r.find(context.Background(), userID, id)
	if err != nil {
		return Automation{}, err
	}
	if len(automations) == 0 {
		return Automation{}, errNotFound
	}
	return automations[0], nil
}

// Create will save new automation into db
func (r *Repo) Create(user string, automation Automation) error {
	return r.create(context.Background(), user, automation)
}

func (r *Repo) CreateTriggerBinding(user string, tb TriggerBinding) error {
	t, err := json.Marshal(tb)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Insert(triggerBindingDbName).
		Cols("userid", "tbid", "triggerid", "info").
		Vals(goqu.Vals{user, tb.TriggerBindingMetadata.ID, tb.TriggerBindingSpec.Trigger, string(t)}).
		ToSQL()
	if err != nil {
		return err
	}

	_, err = r.DB.ExecContext(context.Background(), q)
	if err != nil {
		return err
	}
	return nil
}

// GetTriggerBindingByTriggerID returns TriggerBinding given Trigger ID
func (r *Repo) GetTriggerBindingByTriggerID(userID string, id string) (TriggerBinding, error) {
	res, err := r.findTriggerBinding(context.Background(), userID, id)
	if err != nil {
		return TriggerBinding{}, err
	}
	if len(res) == 0 {
		return TriggerBinding{}, errTBNotFound
	}
	return res[0], nil
}

func (r *Repo) create(ctx context.Context, user string, automation Automation) error {
	t, err := json.Marshal(automation)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Insert(dbName).
		Cols("userid", "automationid", "info").
		Vals(goqu.Vals{user, automation.Metadata.ID, string(t)}).
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
func (r *Repo) find(ctx context.Context, user string, ids ...string) ([]Automation, error) {
	e := exp.Ex{
		"userid": user,
	}
	if len(ids) > 0 {
		e = exp.Ex{
			"userid":       user,
			"automationid": ids,
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
	set := []Automation{}
	for rows.Next() {
		uid := ""
		pid := ""
		t := ""
		if err := rows.Scan(&uid, &pid, &t); err != nil {
			return nil, err
		}
		automation := Automation{}
		json.Unmarshal([]byte(t), &automation)
		set = append(set, automation)
	}
	return set, rows.Close()
}

func (r *Repo) findTriggerBinding(ctx context.Context, user string, ids ...string) ([]TriggerBinding, error) {
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
		From(triggerBindingDbName).
		Where(e).
		ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	set := []TriggerBinding{}
	for rows.Next() {
		uid := ""
		pid := ""
		tid := ""
		t := ""
		if err := rows.Scan(&uid, &pid, &tid, &t); err != nil {
			return nil, err
		}
		tb := TriggerBinding{}
		json.Unmarshal([]byte(t), &tb)
		set = append(set, tb)
	}
	return set, rows.Close()
}
