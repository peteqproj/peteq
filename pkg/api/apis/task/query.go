package task

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/pkg/repo"
)

type (
	// QueryAPI for tasks
	QueryAPI struct {
		Repo *repo.Repo
	}
)

// List returns a list of tasks
// @description Task
// @tags RestAPI
// @produce  json
// @success 200 {array} task.Task
// @router /api/task/ [get]
// @Security ApiKeyAuth
func (a *QueryAPI) List(c *gin.Context) {
	res, err := a.Repo.List(c.Request.Context(), repo.ListOptions{})
	if err != nil {
		handleError(500, err, c)
		return
	}
	c.JSON(200, res)
}

// Get returns a one task
// @description Task
// @tags RestAPI
// @produce  json
// @Param id path string true "Task ID"
// @success 200 {object} task.Task
// @router /api/task/{id} [get]
// @Security ApiKeyAuth
func (a *QueryAPI) Get(c *gin.Context) {
	t, err := a.Repo.Get(c.Request.Context(), repo.GetOptions{ID: c.Param("id")})
	if err != nil {
		handleError(404, fmt.Errorf("Task not found"), c)
		return
	}
	c.JSON(200, t)
}
