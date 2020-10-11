package task

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/peteqproj/peteq/pkg/db/local"
	"github.com/peteqproj/peteq/pkg/logger"
	"gopkg.in/yaml.v2"
)

type (
	// Repo is task repository
	// it works on the view db to read/write from it
	Repo struct {
		DB     *local.DB
		Logger logger.Logger
	}

	// ListOptions to get task list
	ListOptions struct {
		UserID string
		noUser bool
	}
)

// List returns list of tasks
func (r *Repo) List(options ListOptions) ([]Task, error) {
	context, err := r.DB.Read()
	if err != nil {
		return nil, err
	}
	all := []Task{}
	if err := yaml.Unmarshal(context, &all); err != nil {
		return nil, err
	}
	if options.noUser {
		return all, nil
	}

	res := []Task{}
	for _, t := range all {
		if t.Tenant.ID == options.UserID {
			res = append(res, t)
		}
	}
	return res, nil
}

// Get returns task given task id
func (r *Repo) Get(userID string, id string) (Task, error) {
	tasks, err := r.List(ListOptions{UserID: userID})
	if err != nil {
		return Task{}, err
	}
	for _, t := range tasks {
		if t.Metadata.ID == id {
			return t, nil
		}
	}
	return Task{}, fmt.Errorf("Task not found")
}

// Create will save new task into db
func (r *Repo) Create(t Task) error {
	allTasks, err := r.List(ListOptions{noUser: true})
	if err != nil {
		return fmt.Errorf("Failed to load tasks: %w", err)
	}
	set := append(allTasks, t)
	bytes, err := yaml.Marshal(set)
	if err != nil {
		return fmt.Errorf("Failed to marshal task: %w", err)
	}
	if err := r.DB.Write(bytes); err != nil {
		return fmt.Errorf("Failed to persist task to read db: %w", err)
	}
	return nil
}

// Delete will remove task from db
func (r *Repo) Delete(userID string, id string) error {
	allTasks, err := r.List(ListOptions{UserID: userID})
	if err != nil {
		return fmt.Errorf("Failed to load tasks: %w", err)
	}
	index := -1
	for i, t := range allTasks {
		if t.Metadata.ID == id {
			index = i
		}
	}
	if index == -1 {
		return fmt.Errorf("Task not found")
	}
	set := append(allTasks[:index], allTasks[index+1:]...)
	bytes, err := yaml.Marshal(set)
	if err != nil {
		return fmt.Errorf("Failed to marshal task: %w", err)
	}
	if err := r.DB.Write(bytes); err != nil {
		return fmt.Errorf("Failed to persist task to read db: %w", err)
	}
	return nil
}

// Update will update given task
func (r *Repo) Update(t Task) error {
	curr, err := r.Get(t.Tenant.ID, t.Metadata.ID)
	if err != nil {
		return fmt.Errorf("Failed to read previous task: %w", err)
	}
	if err := mergo.Merge(&curr, t, mergo.WithOverwriteWithEmptyValue); err != nil {
		return fmt.Errorf("Failed to update task: %w", err)
	}
	tasks, err := r.List(ListOptions{noUser: true})
	if err != nil {
		return fmt.Errorf("Failed to read tasks: %w", err)
	}
	index := -1
	for i, task := range tasks {
		if task.Metadata.ID == t.Metadata.ID {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("Task not found")
	}
	tasks[index] = curr
	bytes, err := yaml.Marshal(tasks)
	if err != nil {
		return fmt.Errorf("Failed to marshal tasks: %w", err)
	}
	if err := r.DB.Write(bytes); err != nil {
		return fmt.Errorf("Failed to write tasks: %w", err)
	}
	return nil

}
