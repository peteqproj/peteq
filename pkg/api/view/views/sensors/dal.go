package sensors

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"gorm.io/gorm"
)

const dbTableName = "view_sensors"

type (
	DAL struct {
		DB *gorm.DB
	}
)

func (d *DAL) create(ctx context.Context, user string, view sensorsView) error {
	b, err := json.Marshal(view)
	if err != nil {
		return err
	}
	q, _, err := goqu.Insert(dbTableName).Cols("userid", "info").Vals(goqu.Vals{user, string(b)}).ToSQL()
	if err != nil {
		return err
	}
	rows, err := d.DB.Raw(q).Rows()
	if err != nil {
		return err
	}
	return rows.Close()
}

func (d *DAL) load(ctx context.Context, user string) (sensorsView, error) {
	q, _, err := goqu.From(dbTableName).Where(exp.Ex{
		"userid": user,
	}).ToSQL()
	if err != nil {
		return sensorsView{}, fmt.Errorf("Failed to build SQL query: %w", err)
	}
	row := d.DB.Raw(q).Row()
	view := ""
	userid := ""
	if err := row.Scan(&userid, &view); err != nil {
		return sensorsView{}, fmt.Errorf("Failed to scan into sensorsView object: %v", err)
	}
	v := sensorsView{}
	if err := json.Unmarshal([]byte(view), &v); err != nil {
		return v, err
	}
	return v, nil
}

func (d *DAL) update(ctx context.Context, user string, view sensorsView) error {
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
	rows, err := d.DB.Raw(q).Rows()
	if err != nil {
		return fmt.Errorf("Failed to update view_home table: %v", err)
	}
	return rows.Close()
}
