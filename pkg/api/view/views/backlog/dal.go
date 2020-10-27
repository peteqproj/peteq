package backlog

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/peteqproj/peteq/pkg/db"
)

const dbTableName = "view_backlog"

type (
	DAL struct {
		DB db.Database
	}
)

func (d *DAL) load(ctx context.Context, user string) (backlogView, error) {
	q, _, err := goqu.From(dbTableName).Where(exp.Ex{
		"userid": user,
	}).ToSQL()
	if err != nil {
		return backlogView{}, fmt.Errorf("Failed to build SQL query: %w", err)
	}
	row := d.DB.QueryRowContext(ctx, q)

	view := ""
	userid := ""
	if err := row.Scan(&userid, &view); err != nil {
		return backlogView{}, fmt.Errorf("Failed to scan into backlogView object: %v", err)
	}
	if err != nil {
		return backlogView{}, err
	}
	current := backlogView{}
	if err := json.Unmarshal([]byte(view), &current); err != nil {
		return backlogView{}, err
	}
	return current, nil
}
func (d *DAL) create(ctx context.Context, user string, view backlogView) error {
	b, err := json.Marshal(view)
	if err != nil {
		return err
	}
	q, _, err := goqu.Insert(dbTableName).Cols("userid", "info").Vals(goqu.Vals{user, string(b)}).ToSQL()
	if err != nil {
		return err
	}
	rows, err := d.DB.QueryContext(ctx, q)
	if err != nil {
		return err
	}
	return rows.Close()
}
func (d *DAL) update(ctx context.Context, user string, view backlogView) error {
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
		return fmt.Errorf("Failed to update db: %v", err)
	}
	return rows.Close()
}
