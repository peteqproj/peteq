package task

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

var ErrNotFound = errors.New("Task not found")
var errNotInitiated = errors.New("Repository was not initialized, make sure to call Initiate function")
var repoDefEmbed = `name: task
tenant: user
root:
  resource: Task
  database:
    name: task
    postgres:
      dbname: tasks
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
aggregates: []
`
var queries = []string{
	"CREATE TABLE IF NOT EXISTS tasks( id text not null,userid text not null,info json not null,PRIMARY KEY (id));",
	"CREATE INDEX IF NOT EXISTS userid ON tasks ( userid);",
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

func (r *Repo) Create(ctx context.Context, resource *Task) error {
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
		Insert("tasks").
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
func (r *Repo) GetById(ctx context.Context, id string) (*Task, error) {
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

	query, _, err := goqu.From("tasks").Where(e).ToSQL()
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
	resource := &Task{}
	// info column must exist
	if err := json.Unmarshal([]byte(table_info), resource); err != nil {
		return nil, err
	}
	return resource, nil
}
func (r *Repo) UpdateTask(ctx context.Context, resource *Task) error {
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
		Update("tasks").
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
		Delete("tasks").
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

func (r *Repo) ListByUserid(ctx context.Context, userid string) ([]*Task, error) {
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

	sql, _, err := goqu.From("tasks").Where(e).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	res := []*Task{}
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
		resource := &Task{}
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
