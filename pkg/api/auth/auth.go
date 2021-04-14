package auth

import (
	"crypto/sha256"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/pkg/tenant"
)

// IsAuthenticated ensure a request has Authorization header
// it also adds the user to the request.Context
func IsAuthenticated(userRepo *user.Repo) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, map[string]interface{}{
				"error": "unauthorized",
			})
			return
		}
		user, err := userRepo.GetByToken(c.Request.Context(), hash(token))
		if err != nil {
			return
		}
		if user == nil {
			UnauthorizedResponse(c)
			return
		}
		ctx := tenant.ContextWithUser(c.Request.Context(), *user)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func UnauthorizedResponse(c *gin.Context) {
	c.JSON(401, gin.H{
		"error": "unauthorized",
	})
}

func hash(s string) string {
	sh := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sh)
}
