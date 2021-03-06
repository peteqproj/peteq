package saga

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/peteqproj/peteq/domain/project"
	projectCommand "github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/domain/task/command"
	"github.com/peteqproj/peteq/internal/errors"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	rssImporter struct {
		Commandbus  commandbus.CommandBus
		TaskRepo    *task.Repo
		ProjectRepo *project.Repo
		Logger      logger.Logger
		Event       event.Event
		IDGenerator utils.IDGenerator
	}

	rssImporterInput struct {
		Project  string `json:"project"`
		URL      string `json:"url"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	rssGUID          string
	rssImporterState struct {
		tasks map[rssGUID]struct {
			taskID            string
			assignedToProject bool
		}
		project string
	}
)

func (r *rssImporter) Run(ctx context.Context) error {
	r.Logger.Info("Running task rss-importer")
	u := tenant.UserFromContext(ctx)
	if u == nil {
		return errors.ErrMissingUserInContext
	}
	list, err := r.ProjectRepo.ListByUserid(ctx, u.Metadata.ID)
	if err != nil {
		return fmt.Errorf("failed to list projects: %w", err)
	}
	input := &rssImporterInput{}

	s, ok := r.Event.Spec.(map[string]interface{})
	if !ok {
		return fmt.Errorf("given input is not a JSON like string")
	}
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, input)
	if err != nil {
		return err
	}
	projectID := ""

	r.Logger.Info("Parsing RSS", "url", input.URL)
	rss, err := gofeed.NewParser().ParseURL(input.URL)
	if err != nil {
		return err
	}
	if input.Project != "" {
		for _, p := range list {
			if p.Metadata.Name == input.Project {
				r.Logger.Info("Project found", "id", p.Metadata.ID, "name", p.Metadata.Name)
				projectID = p.Metadata.ID
			}
		}
		if projectID == "" {
			id, err := r.IDGenerator.GenerateV4()
			if err != nil {
				return err
			}
			r.Logger.Info("Project was not found, creating", "name", input.Project)
			opt := projectCommand.CreateProjectCommandOptions{
				ID:          id,
				Name:        input.Project,
				Description: "",
				ImageURL:    rss.Image.URL,
			}
			if err := r.Commandbus.Execute(ctx, "project.create", opt); err != nil {
				return fmt.Errorf("Failed to run create.project command: %w", err)
			}
		}
		for _, i := range rss.Items {
			id, err := r.IDGenerator.GenerateV4()
			if err != nil {
				return err
			}
			if err := r.Commandbus.Execute(ctx, "task.create", command.CreateCommandOptions{
				ID:          id,
				Name:        i.Title,
				Description: fmt.Sprintf("Link %s\n", i.Link),
			}); err != nil {
				return err
			}

			if err := r.Commandbus.Execute(ctx, "project.add-task", projectCommand.AddTasksCommandOptions{
				Project: projectID,
				TaskID:  id,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}
