package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/peteqproj/peteq/pkg/api"
	"github.com/peteqproj/peteq/pkg/config"
)

type (
	// Server is http server
	Server interface {
		Start() error
		AddResource(r api.Resource) error
		AddWS(ws *socketio.Server) error
	}

	// Options to create server
	Options struct {
		Config *config.Server
	}

	server struct {
		cnf *config.Server
		srv *gin.Engine
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
	return &server{
		srv: srv,
		cnf: options.Config,
	}
}

func (s *server) Start() error {
	s.srv.Run(fmt.Sprintf("0.0.0.0:%s", s.cnf.Port))
	return nil
}

func (s *server) AddResource(r api.Resource) error {
	return s.addResource(r, nil)
}

func (s *server) AddWS(wsserver *socketio.Server) error {
	s.srv.GET("/ws/*any", gin.WrapH(wsserver))
	s.srv.POST("/ws/*any", gin.WrapH(wsserver))
	return nil
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
