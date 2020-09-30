package project

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// QueryAPI for projects
	QueryAPI struct {
		Repo *project.Repo
	}
)

// List projects
func (q *QueryAPI) List(c *gin.Context) {
	u := tenant.UserFromContext(c.Request.Context())

	res, err := q.Repo.List(project.QueryOptions{
		UserID: u.Metadata.ID,
	})
	if err != nil {
		handleError(500, err, c)
		return
	}
	c.JSON(200, res)
}

// Get returns a one project
func (q *QueryAPI) Get(c *gin.Context) {
	u := tenant.UserFromContext(c.Request.Context())
	p, err := q.Repo.Get(u.Metadata.ID, c.Param("id"))
	if err != nil {
		handleError(404, fmt.Errorf("Project not found"), c)
		return
	}
	c.JSON(200, p)
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}
