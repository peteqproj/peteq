package automation

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/peteqproj/peteq/domain/user"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/logger"
	repo "github.com/peteqproj/peteq/pkg/repo/def"
	"gopkg.in/yaml.v2"

	"github.com/peteqproj/peteq/pkg/tenant"
)

var ErrNotFound = errors.New("Automation not found")
var errNotInitiated = errors.New("Repository was not initialized, make sure to call Initiate function")
var errNoTenantInContext = errors.New("No tenant in context")
var repoDefEmbed = `name: automation
tenant: user
root:
  resource: Automation
  database:
    name: automation_repo
    postgres:
      columns:
      - name: id
        type: text
        fromResource:
          as: string
          path: Metadata.ID
      - name: userid
        type: text
        fromTenant:
          as: string
          path: Metadata.ID
      - name: info
        type: json
        fromResource:
          as: string
          path: .
      indexes:
      - - userid
      uniqueIndexes: []
      primeryKey:
      - id
aggregates:
- resource: TriggerBinding
  database:
    name: trigger_binding_repo
    postgres:
      columns:
      - name: id
        type: text
        fromResource:
          as: string
          path: Metadata.ID
      - name: userid
        type: text
        fromTenant:
          as: string
          path: Metadata.ID
      - name: automation
        type: text
        fromResource:
          as: string
          path: Spec.Automation
      - name: trigger
        type: text
        fromResource:
          as: string
          path: Spec.Trigger
      - name: info
        type: json
        fromResource:
          as: json
          path: .
      indexes:
      - - userid
      uniqueIndexes:
      - - userid
        - trigger
      primeryKey:
      - id
`
var queries = []string{
	"CREATE TABLE IF NOT EXISTS automation_repo( id text not null,userid text not null,info json not null,PRIMARY KEY (id));",
	"CREATE INDEX IF NOT EXISTS userid ON automation_repo ( userid);",
	"CREATE TABLE IF NOT EXISTS trigger_binding_repo( id text not null,userid text not null,automation text not null,trigger text not null,info json not null,PRIMARY KEY (id));",
	"CREATE UNIQUE INDEX IF NOT EXISTS userid_trigger ON trigger_binding_repo ( userid,trigger);",
	"CREATE INDEX IF NOT EXISTS userid ON trigger_binding_repo ( userid);",
}

type (
	Repo struct {
		DB        db.Database
		Logger    logger.Logger
		initiated bool
		def       *repo.RepoDef
	}
)

func (r *Repo) Initiate(ctx context.Context) error {
	for _, q := range queries {
		r.Logger.Info("Running db init query", "query", q)
		if _, err := r.DB.ExecContext(ctx, q); err != nil {
			return err
		}
	}

	def := &repo.RepoDef{}
	if err := yaml.Unmarshal([]byte(repoDefEmbed), def); err != nil {
		return err
	}
	r.def = def

	r.initiated = true
	return nil
}

/* PrimeryKey functions*/

func (r *Repo) Create(ctx context.Context, resource *Automation) error {
	if !r.initiated {
		return errNotInitiated
	}
	var user *user.User
	if r.def.Tenant != "" {
		u := tenant.UserFromContext(ctx)
		if u == nil {
			return errNoTenantInContext
		}
		user = u
	}

	table_column_id := resource.Metadata.ID
	table_column_userid := user.Metadata.ID
	table_column_info, err := json.Marshal(resource)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Insert("automation_repo").
		Cols(
			"id",
			"userid",
			"info",
		).
		Vals(goqu.Vals{
			string(table_column_id),
			string(table_column_userid),
			string(table_column_info),
		}).
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
func (r *Repo) GetById(ctx context.Context, id string) (*Automation, error) {
	if !r.initiated {
		return nil, errNotInitiated
	}
	e := exp.Ex{}
	e["id"] = id
	if r.def.Tenant != "" {
		u := tenant.UserFromContext(ctx)
		if u == nil {
			return nil, errNoTenantInContext
		}
		e["userid"] = u.Metadata.ID
	}

	query, _, err := goqu.From("automation_repo").Where(e).ToSQL()
	if err != nil {
		return nil, err
	}
	row := r.DB.QueryRowContext(ctx, query)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var table_id string
	var table_userid string
	var table_info string

	if err := row.Scan(
		&table_id,
		&table_userid,
		&table_info,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	resource := &Automation{}
	// info column must exist
	if err := json.Unmarshal([]byte(table_info), resource); err != nil {
		return nil, err
	}
	return resource, nil
}
func (r *Repo) UpdateAutomation(ctx context.Context, resource *Automation) error {
	if !r.initiated {
		return errNotInitiated
	}
	var user *user.User
	if r.def.Tenant != "" {
		u := tenant.UserFromContext(ctx)
		if u == nil {
			return errNoTenantInContext
		}
		user = u
	}

	table_column_id := resource.Metadata.ID
	table_column_userid := user.Metadata.ID
	table_column_info, err := json.Marshal(resource)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Update("automation_repo").
		Where(exp.Ex{
			"id": resource.Metadata.ID,
		}).
		Set(goqu.Record{
			"id":     string(table_column_id),
			"userid": string(table_column_userid),
			"info":   string(table_column_info),
		}).
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
func (r *Repo) DeleteById(ctx context.Context, id string) error {
	if !r.initiated {
		return errNotInitiated
	}
	e := exp.Ex{}
	e["id"] = id
	if r.def.Tenant != "" {
		u := tenant.UserFromContext(ctx)
		if u == nil {
			return errNoTenantInContext
		}
		e["userid"] = u.Metadata.ID
	}

	q, _, err := goqu.
		Delete("automation_repo").
		Where(e).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, q)
	return err
}

/*End of PrimeryKey functions*/

/*Index functions*/

func (r *Repo) ListByUserid(ctx context.Context, userid string) ([]*Automation, error) {
	if !r.initiated {
		return nil, errNotInitiated
	}
	e := exp.Ex{}
	e["userid"] = userid
	if r.def.Tenant != "" {
		u := tenant.UserFromContext(ctx)
		if u == nil {
			return nil, errNoTenantInContext
		}
		e["userid"] = u.Metadata.ID
	}

	sql, _, err := goqu.From("automation_repo").Where(e).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	res := []*Automation{}
	for rows.Next() {
		var table_id string
		var table_userid string
		var table_info string

		if err := rows.Scan(
			&table_id,
			&table_userid,
			&table_info,
		); err != nil {
			return nil, err
		}
		resource := &Automation{}
		// info column must exist
		if err := json.Unmarshal([]byte(table_info), resource); err != nil {
			return nil, err
		}
		res = append(res, resource)
	}
	return res, rows.Close()

}

/*End of index function'*/

/*UniqueIndex functions*/
/*End of UniqueIndex functions*/

func (r *Repo) CreateTriggerBinding(ctx context.Context, resource *TriggerBinding) error {
	if !r.initiated {
		return errNotInitiated
	}
	var user *user.User
	if r.def.Tenant != "" {
		u := tenant.UserFromContext(ctx)
		if u == nil {
			return errNoTenantInContext
		}
		user = u
	}

	table_column_id := resource.Metadata.ID
	table_column_userid := user.Metadata.ID
	table_column_automation := resource.Spec.Automation
	table_column_trigger := resource.Spec.Trigger
	table_column_info, err := json.Marshal(resource)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Insert("trigger_binding_repo").
		Cols(
			"id",
			"userid",
			"automation",
			"trigger",
			"info",
		).
		Vals(goqu.Vals{
			string(table_column_id),
			string(table_column_userid),
			string(table_column_automation),
			string(table_column_trigger),
			string(table_column_info),
		}).
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
func (r *Repo) UpdateTriggerBinding(ctx context.Context, resouce *TriggerBinding) error {
	// TODO
	return nil
}

func (r *Repo) GetTriggerBindingById(ctx context.Context, id string) (*TriggerBinding, error) {
	return nil, nil
}

func (r *Repo) ListTriggerBindingByUserid(ctx context.Context, userid string) ([]*TriggerBinding, error) {
	if !r.initiated {
		return nil, errNotInitiated
	}
	e := exp.Ex{}
	e["userid"] = userid
	if r.def.Tenant != "" {
		u := tenant.UserFromContext(ctx)
		if u == nil {
			return nil, errNoTenantInContext
		}
		e["userid"] = u.Metadata.ID
	}

	sql, _, err := goqu.From("trigger_binding_repo").Where(e).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	res := []*TriggerBinding{}
	for rows.Next() {
		var table_id string
		var table_userid string
		var table_automation string
		var table_trigger string
		var table_info string

		if err := rows.Scan(
			&table_id,
			&table_userid,
			&table_automation,
			&table_trigger,
			&table_info,
		); err != nil {
			return nil, err
		}
		resource := &TriggerBinding{}
		// info column must exist
		if err := json.Unmarshal([]byte(table_info), resource); err != nil {
			return nil, err
		}
		res = append(res, resource)
	}
	return res, rows.Close()
}

func (r *Repo) GetTriggerBindingByUseridTrigger(ctx context.Context, userid string, trigger string) (*TriggerBinding, error) {
	if !r.initiated {
		return nil, errNotInitiated
	}
	e := exp.Ex{}
	e["userid"] = userid

	e["trigger"] = trigger
	if r.def.Tenant != "" {
		u := tenant.UserFromContext(ctx)
		if u == nil {
			return nil, errNoTenantInContext
		}
		e["userid"] = u.Metadata.ID
	}

	query, _, err := goqu.From("trigger_binding_repo").Where(e).ToSQL()
	if err != nil {
		return nil, err
	}
	row := r.DB.QueryRowContext(ctx, query)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var table_id string
	var table_userid string
	var table_automation string
	var table_trigger string
	var table_info string

	if err := row.Scan(
		&table_id,
		&table_userid,
		&table_automation,
		&table_trigger,
		&table_info,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	resource := &TriggerBinding{}
	// info column must exist
	if err := json.Unmarshal([]byte(table_info), resource); err != nil {
		return nil, err
	}
	return resource, nil
}
