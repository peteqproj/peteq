package task

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/task"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// QueryAPI for tasks
	QueryAPI struct {
		Repo *task.Repo
	}
)

// List returns a list of tasks
func (a *QueryAPI) List(c *gin.Context) {
	u := tenant.UserFromContext(c.Request.Context())
	res, err := a.Repo.List(task.ListOptions{UserID: u.Metadata.ID})
	if err != nil {
		handleError(500, err, c)
		return
	}
	c.JSON(200, res)
}

// Get returns a one task
func (a *QueryAPI) Get(c *gin.Context) {
	u := tenant.UserFromContext(c.Request.Context())
	t, err := a.Repo.Get(u.Metadata.ID, c.Param("id"))
	if err != nil {
		handleError(404, fmt.Errorf("Task not found"), c)
		return
	}
	c.JSON(200, t)
}
