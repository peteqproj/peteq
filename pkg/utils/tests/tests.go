package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
)

// AuthenticatedContext creates authenticated context with user
func AuthenticatedContext() context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, tenant.User, user.User{
		Metadata: user.Metadata{
			Email: "some@email.com",
			ID:    "user-id",
		},
	})
}

// JSONStringToReadCloser converts json string into io.ReadCloser
// should be used only in test as this method will exit on error
func JSONStringToReadCloser(j map[string]interface{}) io.ReadCloser {
	b, err := json.Marshal(j)
	utils.DieOnError(err, "Failed to convert json to io.ReadCloser")
	return ioutil.NopCloser(bytes.NewReader(b))
}

// MustMarshal marshals or dies
func MustMarshal(v interface{}) []byte {
	r, err := json.Marshal(v)
	if err != nil {
		utils.DieOnError(err, "Failed to marshal")
	}
	return r
}
