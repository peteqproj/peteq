package user

import (
	"context"
	"encoding/json"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

// List selects all the users from the database
func (r *Repo) List(ctx context.Context) ([]*User, error) {
	query, _, err := goqu.From("user_repo").Where(exp.Ex{}).ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.Raw(query).Rows()
	if err != nil {
		return nil, err
	}
	res := []*User{}
	for rows.Next() {
		var table_id string
		var table_email string
		var table_token string
		var table_info string

		if err := rows.Scan(
			&table_id,
			&table_email,
			&table_token,
			&table_info,
		); err != nil {
			return nil, err
		}
		resource := &User{}
		if err := json.Unmarshal([]byte(table_info), resource); err != nil {
			return nil, err
		}
		res = append(res, resource)
	}
	return res, rows.Close()
}
