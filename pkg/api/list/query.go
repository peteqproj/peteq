package list

import (
	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/list"
)

type (
	// QueryAPI for lists
	QueryAPI struct {
		Repo *list.Repo
	}
)

// List lists
func (q *QueryAPI) List(c *gin.Context) {
	res, err := q.Repo.List(list.QueryOptions{})
	if err != nil {
		handleError(500, err, c)
		return
	}
	c.JSON(200, res)
}

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}
