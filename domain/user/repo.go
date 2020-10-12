package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/imdario/mergo"
	"github.com/peteqproj/peteq/pkg/logger"
)

const dbName = "repo_users"

var errNotFound = errors.New("User not found")

type (
	// Repo is user repository
	// it works on the view db to read/write from it
	Repo struct {
		DB     *sql.DB
		Logger logger.Logger
	}

	// ListOptions to get user list
	ListOptions struct{}
)

// List returns list of users
func (r *Repo) List(options ListOptions) ([]User, error) {
	return r.find(context.Background())
}

// Get returns user given user id
func (r *Repo) Get(id string) (User, error) {
	users, err := r.find(context.Background(), id)
	if err != nil {
		return User{}, err
	}
	if len(users) == 0 {
		return User{}, errNotFound
	}
	return users[0], nil
}

// Create will save new user into db
func (r *Repo) Create(u User) error {
	return r.create(context.Background(), u)
}

// Delete will remove user from db
func (r *Repo) Delete(id string) error {
	return r.delete(context.Background(), id)
}

// Update will update given user
func (r *Repo) Update(newUser User) error {
	user, err := r.Get(newUser.Metadata.ID)
	if err != nil {
		return fmt.Errorf("Failed to read previous user: %w", err)
	}
	if err := mergo.Merge(&user, newUser, mergo.WithOverwriteWithEmptyValue); err != nil {
		return fmt.Errorf("Failed to update user: %w", err)
	}
	return r.update(context.Background(), user)
}

func (r *Repo) create(ctx context.Context, user User) error {
	u, err := json.Marshal(user)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Insert(dbName).
		Cols("userid", "info").
		Vals(goqu.Vals{user.Metadata.ID, string(u)}).
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
func (r *Repo) find(ctx context.Context, user ...string) ([]User, error) {
	e := exp.Ex{
		"userid": user,
	}
	if len(user) == 0 {
		e = exp.Ex{}
	}
	q, _, err := goqu.
		From(dbName).
		Where(e).
		ToSQL()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	set := []User{}
	for rows.Next() {
		uid := ""
		u := ""
		if err := rows.Scan(&uid, &u); err != nil {
			return nil, err
		}
		usr := User{}
		json.Unmarshal([]byte(u), &usr)
		set = append(set, usr)
	}
	return set, nil
}
func (r *Repo) delete(ctx context.Context, user string) error {
	q, _, err := goqu.
		Delete(dbName).
		Where(exp.Ex{
			"userid": user,
		}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, q)
	return err
}
func (r *Repo) update(ctx context.Context, user User) error {
	u, err := json.Marshal(user)
	if err != nil {
		return err
	}
	q, _, err := goqu.
		Update(dbName).
		Where(exp.Ex{
			"userid": user.Metadata.ID,
		}).
		Set(goqu.Record{"info": string(u)}).
		ToSQL()
	if err != nil {
		return err
	}
	_, err = r.DB.ExecContext(ctx, q)
	return err
}
