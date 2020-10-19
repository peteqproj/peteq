package project

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

const dbTableName = "view_project"

type (
	DAL struct {
		DB *sql.DB
	}
)

func (d *DAL) create(ctx context.Context, user string, view projectView) error {
	b, err := json.Marshal(view)
	if err != nil {
		return err
	}
	q, _, err := goqu.Insert(dbTableName).Cols("userid", "projectid", "info").Vals(goqu.Vals{user, view.Project.Metadata.ID, string(b)}).ToSQL()
	if err != nil {
		return err
	}
	rows, err := d.DB.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	return rows.Close()
}

func (d *DAL) load(ctx context.Context, user string, project string) (projectView, error) {
	q, _, err := goqu.From(dbTableName).Where(exp.Ex{
		"userid":    user,
		"projectid": project,
	}).ToSQL()
	if err != nil {
		return projectView{}, fmt.Errorf("Failed to build SQL query: %w", err)
	}
	row := d.DB.QueryRowContext(ctx, q)
	view := ""
	userid := ""
	projectid := ""
	if err := row.Scan(&userid, &projectid, &view); err != nil {
		return projectView{}, fmt.Errorf("Failed to scan into projectView object: %v", err)
	}
	v := projectView{}
	if err := json.Unmarshal([]byte(view), &v); err != nil {
		return v, err
	}
	return v, nil
}

func (d *DAL) update(ctx context.Context, user string, project string, view projectView) error {
	res, err := json.Marshal(view)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Update(dbTableName).
		Set(goqu.Record{"info": string(res)}).
		Where(exp.Ex{
			"userid":    user,
			"projectid": project,
		}).
		ToSQL()
	if err != nil {
		return err
	}
	rows, err := d.DB.QueryContext(ctx, q)
	if err != nil {
		return fmt.Errorf("Failed to update view_project table: %v", err)
	}
	return rows.Close()
}
