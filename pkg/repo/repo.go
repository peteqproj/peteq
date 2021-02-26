package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
)

var errNotFound = errors.New("Resource not foud")
var errUserIsMissing = errors.New("User is missing in context")

type (
	// Repo is a connection to repository database for single resource
	Repo struct {
		dbname string
		db     db.Database
		logger logger.Logger
	}

	// ListOptions is all the avaliable options to list resources
	ListOptions struct{}

	listOptions struct {
		ids []string
	}

	// GetOptions is all the avaliable options to get one resource
	GetOptions struct {
		ID string
	}

	// Options to build new repo
	Options struct {
		ResourceType string
		DB           db.Database
		Logger       logger.Logger
	}
)

// New builds repo for given resource
func New(opt Options) (*Repo, error) {
	r := &Repo{
		dbname: fmt.Sprintf("repo_%s", opt.ResourceType),
		db:     opt.DB,
		logger: opt.Logger,
	}
	return r, nil
}

// List all the resources that matches the request
func (r *Repo) List(ctx context.Context, opt ListOptions) ([]*Resource, error) {
	return r.list(ctx, listOptions{})
}

// Get returns one resource that matches the GetOptions request
func (r *Repo) Get(ctx context.Context, opt GetOptions) (*Resource, error) {
	list, err := r.list(ctx, listOptions{
		ids: []string{opt.ID},
	})
	if err != nil {
		return nil, err
	}
	if list == nil || len(list) == 0 {
		return nil, errNotFound
	}
	return list[0], nil
}

// Update overwrites the current resource
func (r *Repo) Update(ctx context.Context, resource Resource) error {
	user := tenant.UserFromContext(ctx)
	if user == nil {
		return errUserIsMissing
	}
	t, err := json.Marshal(resource)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Update(r.dbname).
		Where(exp.Ex{
			"userid":     user.Metadata.ID,
			"resourceid": resource.Metadata.ID,
		}).
		Set(goqu.Record{"info": string(t)}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, q)
	return err
}

// Delete removes the resource
func (r *Repo) Delete(ctx context.Context, resource Resource) error {
	user := tenant.UserFromContext(ctx)
	if user == nil {
		return errUserIsMissing
	}
	q, _, err := goqu.
		Delete(r.dbname).
		Where(exp.Ex{
			"userid":     user.Metadata.ID,
			"resourceid": resource.Metadata.ID,
		}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, q)
	return err
}

// Create stores the resource
func (r *Repo) Create(ctx context.Context, resource Resource) error {
	user := tenant.UserFromContext(ctx)
	if user == nil {
		return errUserIsMissing
	}
	t, err := json.Marshal(resource)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Insert(r.dbname).
		Cols("userid", "resourceid", "info").
		Vals(goqu.Vals{user.Metadata.ID, resource.Metadata.ID, string(t)}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, q)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) list(ctx context.Context, opt listOptions) ([]*Resource, error) {
	user := tenant.UserFromContext(ctx)
	if user == nil {
		return nil, errUserIsMissing
	}
	e := exp.Ex{
		"userid": user.Metadata.ID,
	}
	if len(opt.ids) > 0 {
		e = exp.Ex{
			"userid":     user.Metadata.ID,
			"resourceid": opt.ids,
		}
	}
	q, _, err := goqu.From(r.dbname).Where(e).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []*Resource{}
	for rows.Next() {
		uid := ""
		rui := ""
		data := ""

		if err := rows.Scan(&uid, &rui, &data); err != nil {
			return nil, err
		}
		resource := &Resource{}
		if err := json.Unmarshal([]byte(data), resource); err != nil {
			return nil, err
		}
		list = append(list, resource)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	return list, nil
}
