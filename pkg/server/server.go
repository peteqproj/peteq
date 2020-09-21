package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"github.com/peteqproj/peteq/pkg/api"
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
		Port string
	}

	server struct {
		srv  *gin.Engine
		port string
	}
)

// New build server
func New(options Options) Server {
	port := "8080"
	if options.Port != "" {
		port = options.Port
	}
	srv := gin.Default()
	srv.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()

	})
	return &server{
		srv:  srv,
		port: port,
	}
}

func (s *server) Start() error {
	s.srv.Run(fmt.Sprintf("0.0.0.0:%s", s.port))
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
