package project

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/peteqproj/peteq/pkg/db/local"
	"gopkg.in/yaml.v2"
)

type (
	// Repo is project repository
	// it works on the view db to read/write from it
	Repo struct {
		DB *local.DB
	}

	// QueryOptions to get project project
	QueryOptions struct{}
)

// List returns set of project
func (r *Repo) List(options QueryOptions) ([]Project, error) {
	context, err := r.DB.Read()
	if err != nil {
		return nil, err
	}
	res := []Project{}
	if err := yaml.Unmarshal(context, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Get returns project given project id
func (r *Repo) Get(id string) (Project, error) {
	lists, err := r.List(QueryOptions{})
	if err != nil {
		return Project{}, err
	}
	for _, l := range lists {
		if l.Metadata.ID == id {
			return l, nil
		}
	}
	return Project{}, fmt.Errorf("Project not found")
}

// Create will save new project into db
func (r *Repo) Create(l Project) error {
	allLists, err := r.List(QueryOptions{})
	if err != nil {
		return fmt.Errorf("Failed to load tasks: %w", err)
	}
	set := append(allLists, l)
	bytes, err := yaml.Marshal(set)
	if err != nil {
		return fmt.Errorf("Failed to marshal project: %w", err)
	}
	if err := r.DB.Write(bytes); err != nil {
		return fmt.Errorf("Failed to persist project to read db: %w", err)
	}
	return nil
}

// Delete will remove project from db
func (r *Repo) Delete(id string) error {
	allLists, err := r.List(QueryOptions{})
	if err != nil {
		return fmt.Errorf("Failed to load tasks: %w", err)
	}
	var index *int
	for i, t := range allLists {
		if t.Metadata.ID == id {
			index = &i
		}
	}
	if index == nil {
		return fmt.Errorf("Project not found")
	}
	set := append(allLists[:*index], allLists[*index+1:]...)
	bytes, err := yaml.Marshal(set)
	if err != nil {
		return fmt.Errorf("Failed to marshal project: %w", err)
	}
	if err := r.DB.Write(bytes); err != nil {
		return fmt.Errorf("Failed to persist project to read db: %w", err)
	}
	return nil
}

// Update will update given project
func (r *Repo) Update(p Project) error {
	curr, err := r.Get(p.Metadata.ID)
	if err != nil {
		return fmt.Errorf("Failed to read previous project: %w", err)
	}
	if err := mergo.Merge(&curr, p, mergo.WithOverwriteWithEmptyValue); err != nil {
		return fmt.Errorf("Failed to update project: %w", err)
	}
	lists, err := r.List(QueryOptions{})
	if err != nil {
		return fmt.Errorf("Failed to read lists: %w", err)
	}
	var index *int
	for i, project := range lists {
		if project.Metadata.ID == p.Metadata.ID {
			index = &i
			break
		}
	}
	if index == nil {
		return fmt.Errorf("Project not found")
	}
	lists[*index] = curr
	return r.updateProjects(lists)
}

// AddTask adds task to project
// TODO: check that task is not assigned to other project
func (r *Repo) AddTask(project string, task string) error {
	proj, err := r.Get(project)
	if err != nil {
		return err
	}
	fmt.Println("Updating project")
	proj.Tasks = append(proj.Tasks, task)
	return r.Update(proj)
}

func (r *Repo) updateProjects(set []Project) error {
	bytes, err := yaml.Marshal(set)
	if err != nil {
		return fmt.Errorf("Failed to marshal tasks: %w", err)
	}
	if err := r.DB.Write(bytes); err != nil {
		return fmt.Errorf("Failed to write tasks: %w", err)
	}
	return nil
}
