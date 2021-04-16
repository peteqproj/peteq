package project

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/pkg/api/auth"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// QueryAPI for projects
	QueryAPI struct {
		Repo *project.Repo
	}
)

// List projects
// @description Project
// @tags RestAPI
// @produce  json
// @success 200 {array} project.Project
// @router /api/project [get]
// @Security ApiKeyAuth
func (q *QueryAPI) List(c *gin.Context) {
	u := tenant.UserFromContext(c.Request.Context())
	if u == nil {
		auth.UnauthorizedResponse(c)
		return
	}
	res, err := q.Repo.ListByUserid(c.Request.Context(), u.Metadata.ID)
	if err != nil {
		handleError(500, err, c)
		return
	}
	c.JSON(200, res)
}

// Get returns a one project
// @description Project
// @tags RestAPI
// @produce  json
// @Param id path string true "Project ID"
// @success 200 {object} project.Project
// @router /api/project/{id} [get]
// @Security ApiKeyAuth
func (q *QueryAPI) Get(c *gin.Context) {
	p, err := q.Repo.GetById(c.Request.Context(), c.Param("id"))
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
