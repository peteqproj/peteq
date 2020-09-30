package user

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/peteqproj/peteq/pkg/db/local"
	"github.com/peteqproj/peteq/pkg/logger"
	"gopkg.in/yaml.v2"
)

type (
	// Repo is user repository
	// it works on the view db to read/write from it
	Repo struct {
		DB     *local.DB
		Logger logger.Logger
	}

	// ListOptions to get user list
	ListOptions struct{}
)

// List returns list of users
func (r *Repo) List(options ListOptions) ([]User, error) {
	context, err := r.DB.Read()
	if err != nil {
		return nil, err
	}
	res := []User{}
	if err := yaml.Unmarshal(context, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Get returns user given user id
func (r *Repo) Get(id string) (User, error) {
	users, err := r.List(ListOptions{})
	if err != nil {
		return User{}, err
	}
	for _, u := range users {
		if u.Metadata.ID == id {
			return u, nil
		}
	}
	return User{}, fmt.Errorf("User not found")
}

// Create will save new user into db
func (r *Repo) Create(u User) error {
	users, err := r.List(ListOptions{})
	if err != nil {
		return fmt.Errorf("Failed to load users: %w", err)
	}
	set := append(users, u)
	bytes, err := yaml.Marshal(set)
	if err != nil {
		return fmt.Errorf("Failed to marshal user: %w", err)
	}
	if err := r.DB.Write(bytes); err != nil {
		return fmt.Errorf("Failed to persist user to read db: %w", err)
	}
	return nil
}

// Delete will remove user from db
func (r *Repo) Delete(id string) error {
	users, err := r.List(ListOptions{})
	if err != nil {
		return fmt.Errorf("Failed to load users: %w", err)
	}
	index := -1
	for i, t := range users {
		if t.Metadata.ID == id {
			index = i
		}
	}
	if index == -1 {
		return fmt.Errorf("User not found")
	}
	set := append(users[:index], users[index+1:]...)
	bytes, err := yaml.Marshal(set)
	if err != nil {
		return fmt.Errorf("Failed to marshal user: %w", err)
	}
	if err := r.DB.Write(bytes); err != nil {
		return fmt.Errorf("Failed to persist user to read db: %w", err)
	}
	return nil
}

// Update will update given user
func (r *Repo) Update(t User) error {
	curr, err := r.Get(t.Metadata.ID)
	if err != nil {
		return fmt.Errorf("Failed to read previous user: %w", err)
	}
	if err := mergo.Merge(&curr, t, mergo.WithOverwriteWithEmptyValue); err != nil {
		return fmt.Errorf("Failed to update user: %w", err)
	}
	users, err := r.List(ListOptions{})
	if err != nil {
		return fmt.Errorf("Failed to read users: %w", err)
	}
	index := -1
	for i, user := range users {
		if user.Metadata.ID == t.Metadata.ID {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("User not found")
	}
	users[index] = curr
	bytes, err := yaml.Marshal(users)
	if err != nil {
		return fmt.Errorf("Failed to marshal users: %w", err)
	}
	if err := r.DB.Write(bytes); err != nil {
		return fmt.Errorf("Failed to write users: %w", err)
	}
	return nil

}
