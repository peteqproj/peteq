package tenant

import (
	"context"

	"github.com/peteqproj/peteq/domain/user"
)

type (
	// Tenant may represent a user
	Tenant struct {
		ID   string `json:"id" yaml:"id"`
		Type string `json:"type" type:"type"`
	}
)

// ContextWithUser adds user to the context
func ContextWithUser(ctx context.Context, u user.User) context.Context {
	return context.WithValue(ctx, User, u)
}

// UserFromContext gets user from context if exists
func UserFromContext(ctx context.Context) *user.User {
	u := ctx.Value(User)
	user, ok := u.(user.User)
	if !ok {
		return nil
	}
	return &user
}
