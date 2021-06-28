package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"gorm.io/gorm"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/peteqproj/peteq/pkg/logger"
	repo "github.com/peteqproj/peteq/pkg/repo/def"
	"gopkg.in/yaml.v2"
)

var ErrNotFound = errors.New("User not found")
var errNotInitiated = errors.New("Repository was not initialized, make sure to call Initiate function")
var repoDefEmbed = `name: user
tenant: ""
root:
  resource: User
  database:
    name: user
    postgres:
      dbname: users
      columns:
      - name: id
        type: text
        fromResource:
          as: string
          path: Metadata.ID
      - name: email
        type: text
        fromResource:
          as: string
          path: Spec.Email
      - name: token
        type: text
        fromResource:
          as: string
          path: Spec.TokenHash
      - name: info
        type: json
        fromResource:
          as: string
          path: .
      indexes: []
      uniqueIndexes:
      - - email
      - - token
      primeryKey:
      - id
aggregates: []
`
var queries = []string{
	"CREATE TABLE IF NOT EXISTS user( id text not null,email text not null,token text not null,info json not null,PRIMARY KEY (id));",
	"CREATE UNIQUE INDEX IF NOT EXISTS email ON user ( email);",
	"CREATE UNIQUE INDEX IF NOT EXISTS token ON user ( token);",
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

func (r *Repo) Create(ctx context.Context, resource *User) error {
	if !r.initiated {
		return errNotInitiated
	}

	table_column_id := resource.Metadata.ID
	table_column_email := resource.Spec.Email
	table_column_token := resource.Spec.TokenHash
	table_column_info, err := json.Marshal(resource)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Insert("users").
		Cols(
			"id",
			"email",
			"token",
			"info",
		).
		Vals(goqu.Vals{
			string(table_column_id),
			string(table_column_email),
			string(table_column_token),
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
func (r *Repo) GetById(ctx context.Context, id string) (*User, error) {
	if !r.initiated {
		return nil, errNotInitiated
	}
	e := exp.Ex{}
	e["id"] = id

	query, _, err := goqu.From("users").Where(e).ToSQL()
	if err != nil {
		return nil, err
	}
	row := r.DB.Raw(query).Row()
	if row.Err() != nil {
		return nil, row.Err()
	}
	var table_id string
	var table_email string
	var table_token string
	var table_info string

	if err := row.Scan(
		&table_id,
		&table_email,
		&table_token,
		&table_info,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	resource := &User{}
	// info column must exist
	if err := json.Unmarshal([]byte(table_info), resource); err != nil {
		return nil, err
	}
	return resource, nil
}
func (r *Repo) UpdateUser(ctx context.Context, resource *User) error {
	if !r.initiated {
		return errNotInitiated
	}

	table_column_id := resource.Metadata.ID
	table_column_email := resource.Spec.Email
	table_column_token := resource.Spec.TokenHash
	table_column_info, err := json.Marshal(resource)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Update("users").
		Where(exp.Ex{
			"id": resource.Metadata.ID,
		}).
		Set(goqu.Record{
			"id":    string(table_column_id),
			"email": string(table_column_email),
			"token": string(table_column_token),
			"info":  string(table_column_info),
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

	q, _, err := goqu.
		Delete("users").
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

/*End of index function'*/

/*UniqueIndex functions*/
func (r *Repo) GetByEmail(ctx context.Context, email string) (*User, error) {
	if !r.initiated {
		return nil, errNotInitiated
	}
	e := exp.Ex{}
	e["email"] = email

	query, _, err := goqu.From("users").Where(e).ToSQL()
	if err != nil {
		return nil, err
	}
	row := r.DB.Raw(query).Row()
	if row.Err() != nil {
		return nil, row.Err()
	}
	var table_id string
	var table_email string
	var table_token string
	var table_info string

	if err := row.Scan(
		&table_id,
		&table_email,
		&table_token,
		&table_info,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	resource := &User{}
	// info column must exist
	if err := json.Unmarshal([]byte(table_info), resource); err != nil {
		return nil, err
	}
	return resource, nil
}
func (r *Repo) GetByToken(ctx context.Context, token string) (*User, error) {
	if !r.initiated {
		return nil, errNotInitiated
	}
	e := exp.Ex{}
	e["token"] = token

	query, _, err := goqu.From("users").Where(e).ToSQL()
	if err != nil {
		return nil, err
	}
	row := r.DB.Raw(query).Row()
	if row.Err() != nil {
		return nil, row.Err()
	}
	var table_id string
	var table_email string
	var table_token string
	var table_info string

	if err := row.Scan(
		&table_id,
		&table_email,
		&table_token,
		&table_info,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	resource := &User{}
	// info column must exist
	if err := json.Unmarshal([]byte(table_info), resource); err != nil {
		return nil, err
	}
	return resource, nil
}

/*End of UniqueIndex functions*/
