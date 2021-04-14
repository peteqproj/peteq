package saga

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mmcdole/gofeed"
	projectCommand "github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/domain/task/command"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/event"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/repo"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	rssImporter struct {
		Commandbus  commandbus.CommandBus
		TaskRepo    *task.Repo
		ProjectRepo *repo.Repo
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
)

func (r *rssImporter) Run(ctx context.Context) error {
	r.Logger.Info("Running task rss-importer")
	usr := tenant.UserFromContext(ctx)
	if usr == nil {
		return fmt.Errorf("User was not set in the current context")
	}
	list, err := r.ProjectRepo.List(ctx, repo.ListOptions{})
	if err != nil {
		return fmt.Errorf("Failed to list projects: %w", err)
	}
	input := &rssImporterInput{}

	s, ok := r.Event.Spec.(string)
	if !ok {
		return fmt.Errorf("Given input is not a JSON like string")
	}
	err = json.Unmarshal([]byte(s), input)
	if err != nil {
		return err
	}
	projectID := ""

	r.Logger.Info("Pasing RSS", "url", input.URL)
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
		if projectID != "" {
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
		}
	}

	return nil
}
