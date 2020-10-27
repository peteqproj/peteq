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
	"github.com/peteqproj/peteq/pkg/api/auth"
	"github.com/peteqproj/peteq/pkg/api/view"
	viewBuilder "github.com/peteqproj/peteq/pkg/api/view/builder"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/db"
	"github.com/peteqproj/peteq/pkg/event/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// Builder builds apis
	Builder struct {
		UserRepo    *user.Repo
		ListRpeo    *list.Repo
		ProjectRepo *project.Repo
		TaskRepo    *task.Repo
		Commandbus  commandbus.CommandBus
		Eventbus    bus.Eventbus
		Logger      logger.Logger
		DB          db.Database
	}
)

// BuildCommandAPI builds command api
func (b *Builder) BuildCommandAPI() api.Resource {
	idGen := utils.NewGenerator()
	taskCommandAPI := taskAPI.CommandAPI{
		Repo:        b.TaskRepo,
		Commandbus:  b.Commandbus,
		Logger:      b.Logger.Fork("api", "task"),
		IDGenerator: idGen,
	}
	listCommandAPI := listAPI.CommandAPI{
		Repo:        b.ListRpeo,
		Commandbus:  b.Commandbus,
		Logger:      b.Logger.Fork("api", "list"),
		IDGenerator: idGen,
	}
	projectCommandAPI := projectAPI.CommandAPI{
		Repo:        b.ProjectRepo,
		Commandbus:  b.Commandbus,
		Logger:      b.Logger.Fork("api", "project"),
		IDGenerator: idGen,
	}
	userCommandAPI := userAPI.CommandAPI{
		Repo:        b.UserRepo,
		Commandbus:  b.Commandbus,
		Logger:      b.Logger.Fork("api", "user"),
		IDGenerator: idGen,
	}
	return api.Resource{
		Path: "/c",
		Subresource: []api.Resource{
			{
				Path: "/task",
				Midderwares: []gin.HandlerFunc{
					auth.IsAuthenticated(b.UserRepo),
				},
				Endpoints: []api.Endpoint{
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
			},
			{
				Path: "/list",
				Midderwares: []gin.HandlerFunc{
					auth.IsAuthenticated(b.UserRepo),
				},
				Endpoints: []api.Endpoint{
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
			},
			{
				Path: "/user",
				Endpoints: []api.Endpoint{
					{
						Path:    "/register",
						Verb:    "POST",
						Handler: api.WrapCommandAPI(userCommandAPI.Register, b.Logger),
					},
					{
						Path:    "/login",
						Verb:    "POST",
						Handler: api.WrapCommandAPI(userCommandAPI.Login, b.Logger),
					},
				},
			},
		},
	}
}

// BuildViewAPI builds view apis
func (b *Builder) BuildViewAPI() api.Resource {
	vb := viewBuilder.New(&viewBuilder.Options{
		TaskRepo:    b.TaskRepo,
		ListRepo:    b.ListRpeo,
		ProjectRepo: b.ProjectRepo,
		Logger:      b.Logger,
		DB:          b.DB,
	})
	views := vb.BuildViews()
	for _, view := range views {
		for name, handler := range view.EventHandlers() {
			b.Logger.Info("Subscribing", "name", name, "handler", handler.Name())
			b.Eventbus.Subscribe(name, handler)
		}
	}
	resource := api.Resource{
		Path: "/q",
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

// BuildRestfulAPI builds restful apis
func (b *Builder) BuildRestfulAPI() api.Resource {
	taskQueryAPI := taskAPI.QueryAPI{
		Repo: b.TaskRepo,
	}
	listQueryAPI := listAPI.QueryAPI{
		Repo: b.ListRpeo,
	}
	projectQueryAPI := projectAPI.QueryAPI{
		Repo: b.ProjectRepo,
	}
	return api.Resource{
		Path: "/api",
		Midderwares: []gin.HandlerFunc{
			auth.IsAuthenticated(b.UserRepo),
		},
		Subresource: []api.Resource{
			{
				Path: "/task",
				Endpoints: []api.Endpoint{
					{
						Verb:    "GET",
						Path:    "/",
						Handler: taskQueryAPI.List,
					},
					{
						Verb:    "GET",
						Path:    "/:id",
						Handler: taskQueryAPI.Get,
					},
				},
			},
			{
				Path: "/project",
				Endpoints: []api.Endpoint{
					{
						Verb:    "GET",
						Handler: projectQueryAPI.List,
					},
					{
						Path:    "/:id",
						Verb:    "GET",
						Handler: projectQueryAPI.Get,
					},
				},
			},
			{
				Path: "/list",
				Endpoints: []api.Endpoint{
					{
						Verb:    "GET",
						Handler: listQueryAPI.List,
					},
				},
			},
		},
	}
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
