package test

import (
	"context"

	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/pkg/tenant"
)

func BuildAuthenticationContextWithUser() context.Context {
	ctx := context.Background()
	return tenant.ContextWithUser(ctx, user.User{
		Metadata: user.Metadata{
			ID: "fake-user-id",
		},
		Spec: user.Spec{
			Email: "test@test.com",
		},
	})
}
