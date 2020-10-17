package api

import "github.com/gin-gonic/gin"

type (
	Resource struct {
		// Path, including /
		Path        string
		Midderwares []gin.HandlerFunc
		Endpoints   []Endpoint
		Subresource []Resource
	}

	// Endpoint is one endpoint
	Endpoint struct {
		// HTTP verb
		Verb string
		// Path, including /
		Path        string
		Midderwares []gin.HandlerFunc
		Handler     gin.HandlerFunc
	}
)
