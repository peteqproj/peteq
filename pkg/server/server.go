package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/peteqproj/peteq/docs"
	"github.com/peteqproj/peteq/pkg/api"
	"github.com/peteqproj/peteq/pkg/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type (
	// Server is http server
	Server interface {
		Start() error
		AddResource(r api.Resource) error
		SetReady()
	}

	// Options to create server
	Options struct {
		Config *config.Server
	}

	server struct {
		cnf     *config.Server
		srv     *gin.Engine
		isReady bool
	}
)

// New build server
func New(options Options) Server {
	srv := gin.Default()
	srv.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()

	})
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	srv.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	s := &server{
		srv:     srv,
		cnf:     options.Config,
		isReady: false,
	}
	s.AddResource(api.Resource{
		Path: "/ready",
		Endpoints: []api.Endpoint{
			{
				Verb: "GET",
				Path: "/",
				Handler: func(c *gin.Context) {
					if s.isReady {
						c.Status(200)
						return
					}
					c.Status(500)
					return
				},
			},
		},
	})

	return s
}

func (s *server) Start() error {
	s.srv.Run(fmt.Sprintf("0.0.0.0:%s", s.cnf.Port))
	return nil
}

func (s *server) AddResource(r api.Resource) error {
	return s.addResource(r, nil)
}

func (s *server) SetReady() {
	s.isReady = true
}

func (s *server) addResource(r api.Resource, parent *gin.RouterGroup) error {
	var router *gin.RouterGroup
	if parent == nil {
		router = s.srv.Group(r.Path, r.Midderwares...)
	} else {
		router = parent.Group(r.Path, r.Midderwares...)
	}
	for _, sg := range r.Subresource {
		if err := s.addResource(sg, router); err != nil {
			return fmt.Errorf("Failed to add resource %s: %w", r.Path, err)
		}
	}

	for _, ep := range r.Endpoints {
		mds := append(ep.Midderwares, ep.Handler)
		router.Handle(ep.Verb, ep.Path, mds...)
	}

	return nil
}
