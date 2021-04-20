package list

import (
	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/pkg/api/auth"
	"github.com/peteqproj/peteq/pkg/tenant"
)

type (
	// QueryAPI for lists
	QueryAPI struct {
		Repo *list.Repo
	}
)

// List lists
// @description List
// @tags RestAPI
// @produce  json
// @success 200 {array} list.List
// @router /api/list [get]
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

func handleError(code int, err error, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": err.Error(),
	})
}
