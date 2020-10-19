package view

import (
	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/pkg/event/handler"
)

type (
	// View can only retrieve data
	View interface {
		Get(c *gin.Context)
		EventHandlers() map[string]handler.EventHandler
	}
)
