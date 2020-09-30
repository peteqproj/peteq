package list

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/peteqproj/peteq/pkg/db/local"
	"github.com/peteqproj/peteq/pkg/logger"
	"gopkg.in/yaml.v2"
)

type (
	// Repo is list repository
	// it works on the view db to read/write from it
	Repo struct {
		DB     *local.DB
		Logger logger.Logger
	}

	// QueryOptions to get task list
	QueryOptions struct {
		UserID string
		noUser bool
	}
)

// List returns set of list
func (r *Repo) List(options QueryOptions) ([]List, error) {
	context, err := r.DB.Read()
	if err != nil {
		return nil, err
	}
	all := []List{}
	if err := yaml.Unmarshal(context, &all); err != nil {
		return nil, err
	}
	res := []List{}
	for _, l := range all {
		if l.Tenant.ID == options.UserID {
			res = append(res, l)
		}
	}
	return res, nil
}

// Get returns list given list id
func (r *Repo) Get(userID string, id string) (List, error) {
	lists, err := r.List(QueryOptions{UserID: userID})
	if err != nil {
		return List{}, err
	}
	for _, l := range lists {
		if l.Metadata.ID == id {
			return l, nil
		}
	}
	return List{}, fmt.Errorf("Task not found")
}

// Create will save new task into db
func (r *Repo) Create(l List) error {
	allLists, err := r.List(QueryOptions{noUser: true})
	if err != nil {
		return fmt.Errorf("Failed to load tasks: %w", err)
	}
	set := append(allLists, l)
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
	allLists, err := r.List(QueryOptions{UserID: userID})
	if err != nil {
		return fmt.Errorf("Failed to load tasks: %w", err)
	}
	index := -1
	for i, t := range allLists {
		if t.Metadata.ID == id {
			index = i
		}
	}
	if index == -1 {
		return fmt.Errorf("Task not found")
	}
	set := append(allLists[:index], allLists[index+1:]...)
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
func (r *Repo) Update(l List) error {
	curr, err := r.Get(l.Tenant.ID, l.Metadata.ID)
	if err != nil {
		return fmt.Errorf("Failed to read previous task: %w", err)
	}
	if err := mergo.Merge(&curr, l, mergo.WithOverwriteWithEmptyValue); err != nil {
		return fmt.Errorf("Failed to update task: %w", err)
	}
	lists, err := r.List(QueryOptions{noUser: true})
	if err != nil {
		return fmt.Errorf("Failed to read lists: %w", err)
	}
	index := -1
	for i, list := range lists {
		if list.Metadata.ID == l.Metadata.ID {
			index = &
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("Task not found")
	}
	lists[index] = curr
	return r.updateLists(lists)
}

// MoveTask will move tasks from one list to another one
// TODO: Validation source and destination are exists
func (r *Repo) MoveTask(userID string, source string, destination string, task string) error {
	listSet, err := r.List(QueryOptions{UserID: userID})
	if err != nil {
		return fmt.Errorf("Failed to load lists: %w", err)
	}
	var sourceListIndex int = -1
	var destinationListIndex int = -1
	for i, l := range listSet {
		if l.Metadata.ID == source {
			sourceListIndex = i
		}

		if l.Metadata.ID == destination {
			destinationListIndex = i
		}
	}

	// If source found, remove task from source
	if sourceListIndex != -1 {
		for i, tid := range listSet[sourceListIndex].Tasks {
			if tid == task {
				listSet[sourceListIndex].Tasks = remove(listSet[sourceListIndex].Tasks, i)
				break
			}
		}
	}

	// If destination found add it to destination
	if destinationListIndex != -1 {
		listSet[destinationListIndex].Tasks = append(listSet[destinationListIndex].Tasks, task)
	}

	err = r.updateLists(listSet)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) updateLists(set []List) error {
	bytes, err := yaml.Marshal(set)
	if err != nil {
		return fmt.Errorf("Failed to marshal tasks: %w", err)
	}
	if err := r.DB.Write(bytes); err != nil {
		return fmt.Errorf("Failed to write tasks: %w", err)
	}
	return nil
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
