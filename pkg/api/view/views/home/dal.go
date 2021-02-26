package home

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/repo"
)

const dbTableName = "view_home"

type (
	DAL struct {
		DB db.Database
	}

	listCreatedHandler struct {
		dal *DAL
	}
	listTaskMovedHandler struct {
		dal         *DAL
		taskRepo    *repo.Repo
		projectRepo *repo.Repo
	}
	taskCreatedHandler struct {
		dal *DAL
	}
	taskUpdateHandler struct {
		dal *DAL
	}
	taskDeletedHandler struct {
		dal      *DAL
		taskRepo *repo.Repo
	}
	userRegistredHandler struct {
		dal *DAL
	}
	projectTaskAddedHandler struct {
		dal         *DAL
		projectRepo *repo.Repo
		taskRepo    *repo.Repo
	}
)

func (d *DAL) create(ctx context.Context, user string, view homeView) error {
	b, err := json.Marshal(view)
	if err != nil {
		return err
	}
	q, _, err := goqu.Insert(dbTableName).Cols("userid", "info").Vals(goqu.Vals{user, string(b)}).ToSQL()
	if err != nil {
		return err
	}
	res, err := d.DB.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	return res.Close()
}
func (d *DAL) load(ctx context.Context, user string) (homeView, error) {
	q, _, err := goqu.From(dbTableName).Where(exp.Ex{
		"userid": user,
	}).ToSQL()
	if err != nil {
		return homeView{}, fmt.Errorf("Failed to build SQL query: %w", err)
	}
	row := d.DB.QueryRowContext(ctx, q)
	view := ""
	userid := ""
	if err := row.Scan(&userid, &view); err != nil {
		return homeView{}, fmt.Errorf("Failed to scan into homeView object: %v", err)
	}
	v := homeView{}
	if err := json.Unmarshal([]byte(view), &v); err != nil {
		return v, err
	}
	return v, nil
}
func (d *DAL) update(ctx context.Context, user string, view homeView) error {
	res, err := json.Marshal(view)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Update(dbTableName).
		Set(goqu.Record{"info": string(res)}).
		Where(exp.Ex{
			"userid": user,
		}).
		ToSQL()
	if err != nil {
		return err
	}
	rows, err := d.DB.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("Failed to update view_home table: %v", err)
	}
	return rows.Close()
}
