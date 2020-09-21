package task

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/task"
)

type (
	// QueryAPI for tasks
	QueryAPI struct {
		Repo *task.Repo
	}
)

// List returns a list of tasks
func (a *QueryAPI) List(c *gin.Context) {
	res, err := a.Repo.List(task.ListOptions{})
	if err != nil {
		handleError(500, err, c)
		return
	}
	c.JSON(200, res)
}

// Get returns a one task
func (a *QueryAPI) Get(c *gin.Context) {
	t, err := a.Repo.Get(c.Param("id"))
	if err != nil {
		handleError(404, fmt.Errorf("Task not found"), c)
		return
	}
	c.JSON(200, t)
}
