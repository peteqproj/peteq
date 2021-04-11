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
		users, err := userRepo.List(user.ListOptions{})
		if err != nil {
			return
		}
		hashed := hash(token)
		for _, u := range users {
			if hashed == u.Spec.TokenHash {
				ctx := tenant.ContextWithUser(c.Request.Context(), u)
				c.Request = c.Request.WithContext(ctx)
				c.Next()
				return
			}
		}
		UnauthorizedResponse(c)
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
