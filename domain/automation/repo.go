package automation

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"gorm.io/gorm"

	"github.com/peteqproj/peteq/domain/user"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	perrors "github.com/peteqproj/peteq/internal/errors"
	"github.com/peteqproj/peteq/pkg/logger"
	repo "github.com/peteqproj/peteq/pkg/repo/def"
	"gopkg.in/yaml.v2"

	"github.com/peteqproj/peteq/pkg/tenant"
)

var ErrNotFound = errors.New("Automation not found")
var errNotInitiated = errors.New("Repository was not initialized, make sure to call Initiate function")
var repoDefEmbed = `name: automation
tenant: user
root:
  resource: Automation
  database:
    name: automation
    postgres:
      dbname: automations
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
- resource: SensorBinding
  database:
    name: sensor_binding
    postgres:
      dbname: sensor_bindings
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
      - name: sensor
        type: text
        fromResource:
          as: string
          path: Spec.Sensor
      - name: info
        type: json
        fromResource:
          as: json
          path: .
      indexes:
      - - userid
      uniqueIndexes:
      - - userid
        - sensor
      primeryKey:
      - id
`
var queries = []string{
	"CREATE TABLE IF NOT EXISTS automations( id text not null,userid text not null,info json not null,PRIMARY KEY (id));",
	"CREATE INDEX IF NOT EXISTS userid ON automations ( userid);",
	"CREATE TABLE IF NOT EXISTS sensor_bindings( id text not null,userid text not null,automation text not null,sensor text not null,info json not null,PRIMARY KEY (id));",
	"CREATE UNIQUE INDEX IF NOT EXISTS userid_sensor ON sensor_bindings ( userid,sensor);",
	"CREATE INDEX IF NOT EXISTS userid ON sensor_bindings ( userid);",
}

type (
	Repo struct {
		DB        *gorm.DB
		Logger    logger.Logger
		initiated bool
		def       *repo.RepoDef
	}
)

func (r *Repo) Initiate(ctx context.Context) error {
	for _, q := range queries {
		r.Logger.Info("Running db init query", "query", q)
		res := r.DB.Exec(q)
		if res.Error != nil {
			return res.Error
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
			return perrors.ErrMissingUserInContext
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
		Insert("automations").
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
	_, err = r.DB.Raw(q).Rows()
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
			return nil, perrors.ErrMissingUserInContext
		}
		e["userid"] = u.Metadata.ID
	}

	query, _, err := goqu.From("automations").Where(e).ToSQL()
	if err != nil {
		return nil, err
	}
	row := r.DB.Raw(query).Row()
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
			return perrors.ErrMissingUserInContext
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
		Update("automations").
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
	_, err = r.DB.Raw(q).Rows()
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
			return perrors.ErrMissingUserInContext
		}
		e["userid"] = u.Metadata.ID
	}

	q, _, err := goqu.
		Delete("automations").
		Where(e).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.DB.Raw(q).Rows()
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
			return nil, perrors.ErrMissingUserInContext
		}
		e["userid"] = u.Metadata.ID
	}

	sql, _, err := goqu.From("automations").Where(e).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.Raw(sql).Rows()
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

func (r *Repo) CreateSensorBinding(ctx context.Context, resource *SensorBinding) error {
	if !r.initiated {
		return errNotInitiated
	}
	var user *user.User
	if r.def.Tenant != "" {
		u := tenant.UserFromContext(ctx)
		if u == nil {
			return perrors.ErrMissingUserInContext
		}
		user = u
	}

	table_column_id := resource.Metadata.ID
	table_column_userid := user.Metadata.ID
	table_column_automation := resource.Spec.Automation
	table_column_sensor := resource.Spec.Sensor
	table_column_info, err := json.Marshal(resource)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Insert("sensor_bindings").
		Cols(
			"id",
			"userid",
			"automation",
			"sensor",
			"info",
		).
		Vals(goqu.Vals{
			string(table_column_id),
			string(table_column_userid),
			string(table_column_automation),
			string(table_column_sensor),
			string(table_column_info),
		}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.DB.Raw(q).Rows()
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) UpdateSensorBinding(ctx context.Context, resouce *SensorBinding) error {
	// TODO
	return nil
}

func (r *Repo) GetSensorBindingById(ctx context.Context, id string) (*SensorBinding, error) {
	return nil, nil
}

func (r *Repo) ListSensorBindingByUserid(ctx context.Context, userid string) ([]*SensorBinding, error) {
	if !r.initiated {
		return nil, errNotInitiated
	}
	e := exp.Ex{}
	e["userid"] = userid
	if r.def.Tenant != "" {
		u := tenant.UserFromContext(ctx)
		if u == nil {
			return nil, perrors.ErrMissingUserInContext
		}
		e["userid"] = u.Metadata.ID
	}

	sql, _, err := goqu.From("sensor_bindings").Where(e).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	res := []*SensorBinding{}
	for rows.Next() {
		var table_id string
		var table_userid string
		var table_automation string
		var table_sensor string
		var table_info string

		if err := rows.Scan(
			&table_id,
			&table_userid,
			&table_automation,
			&table_sensor,
			&table_info,
		); err != nil {
			return nil, err
		}
		resource := &SensorBinding{}
		// info column must exist
		if err := json.Unmarshal([]byte(table_info), resource); err != nil {
			return nil, err
		}
		res = append(res, resource)
	}
	return res, rows.Close()
}

func (r *Repo) GetSensorBindingByUseridSensor(ctx context.Context, userid string, sensor string) (*SensorBinding, error) {
	if !r.initiated {
		return nil, errNotInitiated
	}
	e := exp.Ex{}
	e["userid"] = userid

	e["sensor"] = sensor
	if r.def.Tenant != "" {
		u := tenant.UserFromContext(ctx)
		if u == nil {
			return nil, perrors.ErrMissingUserInContext
		}
		e["userid"] = u.Metadata.ID
	}

	query, _, err := goqu.From("sensor_bindings").Where(e).ToSQL()
	if err != nil {
		return nil, err
	}
	row := r.DB.Raw(query).Row()
	if row.Err() != nil {
		return nil, row.Err()
	}
	var table_id string
	var table_userid string
	var table_automation string
	var table_sensor string
	var table_info string

	if err := row.Scan(
		&table_id,
		&table_userid,
		&table_automation,
		&table_sensor,
		&table_info,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	resource := &SensorBinding{}
	// info column must exist
	if err := json.Unmarshal([]byte(table_info), resource); err != nil {
		return nil, err
	}
	return resource, nil
}
