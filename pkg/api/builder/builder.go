package builder

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/pkg/api"
	listAPI "github.com/peteqproj/peteq/pkg/api/apis/list"
	projectAPI "github.com/peteqproj/peteq/pkg/api/apis/project"
	taskAPI "github.com/peteqproj/peteq/pkg/api/apis/task"
	userAPI "github.com/peteqproj/peteq/pkg/api/apis/user"
	"github.com/peteqproj/peteq/pkg/api/apis/view"
	"github.com/peteqproj/peteq/pkg/api/auth"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// Builder builds apis
	Builder struct {
		UserRepo    *user.Repo
		ListRpeo    *list.Repo
		ProjectRepo *project.Repo
		TaskRepo    *task.Repo
		Commandbus  commandbus.CommandBus
		Logger      logger.Logger
	}
)

// BuildCommandAPI builds restful apis
func (b *Builder) BuildCommandAPI() api.Resource {
	taskCommandAPI := taskAPI.CommandAPI{
		Repo:       b.TaskRepo,
		Commandbus: b.Commandbus,
		Logger:     b.Logger.Fork("api", "task"),
	}
	taskQueryAPI := taskAPI.QueryAPI{
		Repo: b.TaskRepo,
	}
	listCommandAPI := listAPI.CommandAPI{
		Repo:       b.ListRpeo,
		Commandbus: b.Commandbus,
		Logger:     b.Logger.Fork("api", "list"),
	}
	listQueryAPI := listAPI.QueryAPI{
		Repo: b.ListRpeo,
	}
	projectCommandAPI := projectAPI.CommandAPI{
		Repo:       b.ProjectRepo,
		Commandbus: b.Commandbus,
		Logger:     b.Logger.Fork("api", "project"),
	}
	projectQueryAPI := projectAPI.QueryAPI{
		Repo: b.ProjectRepo,
	}
	userCommandAPI := userAPI.CommandAPI{
		Repo:       b.UserRepo,
		Commandbus: b.Commandbus,
		Logger:     b.Logger.Fork("api", "user"),
	}
	return api.Resource{
		Path: "/api",
		Subresource: []api.Resource{
			{
				Path: "/task",
				Midderwares: []gin.HandlerFunc{
					auth.IsAuthenticated(b.UserRepo),
				},
				Endpoints: []api.Endpoint{
					{
						Verb:    "GET",
						Handler: taskQueryAPI.List,
					},
					{
						Verb:    "POST",
						Path:    "/complete",
						Handler: api.WrapCommandAPI(taskCommandAPI.Complete, b.Logger),
					},
					{
						Verb:    "POST",
						Path:    "/reopen",
						Handler: api.WrapCommandAPI(taskCommandAPI.Reopen, b.Logger),
					},
					{
						Verb:    "POST",
						Path:    "/create",
						Handler: api.WrapCommandAPI(taskCommandAPI.Create, b.Logger),
					},
					{
						Verb:    "POST",
						Path:    "/update",
						Handler: api.WrapCommandAPI(taskCommandAPI.Update, b.Logger),
					},
					{
						Verb:    "POST",
						Path:    "/delete",
						Handler: api.WrapCommandAPI(taskCommandAPI.Delete, b.Logger),
					},
				},
				Subresource: []api.Resource{
					{
						Path: "/:id",
						Endpoints: []api.Endpoint{
							{
								Verb:    "GET",
								Handler: taskQueryAPI.Get,
							},
						},
					},
				},
			},
			{
				Path: "/list",
				Midderwares: []gin.HandlerFunc{
					auth.IsAuthenticated(b.UserRepo),
				},
				Endpoints: []api.Endpoint{
					{
						Verb:    "GET",
						Handler: listQueryAPI.List,
					},
					{
						Verb:    "POST",
						Path:    "/moveTasks",
						Handler: api.WrapCommandAPI(listCommandAPI.MoveTasks, b.Logger),
					},
				},
			},
			{
				Path: "/project",
				Midderwares: []gin.HandlerFunc{
					auth.IsAuthenticated(b.UserRepo),
				},
				Endpoints: []api.Endpoint{
					{
						Verb:    "GET",
						Handler: projectQueryAPI.List,
					},
					{
						Path:    "/create",
						Verb:    "POST",
						Handler: api.WrapCommandAPI(projectCommandAPI.Create, b.Logger),
					},
					{
						Path:    "/addTasks",
						Verb:    "POST",
						Handler: api.WrapCommandAPI(projectCommandAPI.AddTasks, b.Logger),
					},
				},
				Subresource: []api.Resource{
					{
						Path: "/:id",
						Endpoints: []api.Endpoint{
							{
								Verb:    "GET",
								Handler: projectQueryAPI.Get,
							},
						},
					},
				},
			},
			{
				Path: "/user",
				Endpoints: []api.Endpoint{
					{
						Path:    "/register",
						Verb:    "POST",
						Handler: api.WrapCommandAPI(userCommandAPI.Register, b.Logger),
					},
				},
			},
		},
	}
}

// BuildViewAPI builds view apis
func (b *Builder) BuildViewAPI() api.Resource {
	views := map[string]view.View{
		"backlog": view.NewView(view.Options{
			Type:        "backlog",
			TaskRepo:    b.TaskRepo,
			ListRepo:    b.ListRpeo,
			ProjectRepo: b.ProjectRepo,
			Logger:      b.Logger,
		}),
		"projects": view.NewView(view.Options{
			Type:        "projects",
			TaskRepo:    b.TaskRepo,
			ListRepo:    b.ListRpeo,
			ProjectRepo: b.ProjectRepo,
			Logger:      b.Logger,
		}),
		"projects/:id": view.NewView(view.Options{
			Type:        "project",
			TaskRepo:    b.TaskRepo,
			ListRepo:    b.ListRpeo,
			ProjectRepo: b.ProjectRepo,
			Logger:      b.Logger,
		}),
		"home": view.NewView(view.Options{
			Type:        "home",
			TaskRepo:    b.TaskRepo,
			ListRepo:    b.ListRpeo,
			ProjectRepo: b.ProjectRepo,
			Logger:      b.Logger,
		}),
	}
	resource := api.Resource{
		Path: "/view",
		Midderwares: []gin.HandlerFunc{
			auth.IsAuthenticated(b.UserRepo),
		},
		Subresource: []api.Resource{
			{
				Endpoints: buildViews(views),
			},
		},
	}
	return resource
}

func buildViews(views map[string]view.View) []api.Endpoint {
	endpoints := []api.Endpoint{}
	for name, view := range views {
		endpoints = append(endpoints, api.Endpoint{
			Path:    fmt.Sprintf("/%s", name),
			Verb:    "GET",
			Handler: view.Get,
		})
	}
	return endpoints
}
