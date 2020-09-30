package user

import (
	"context"
	"crypto/sha256"
	"io"

	"github.com/gofrs/uuid"
	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
)

type (
	// CommandAPI for users
	CommandAPI struct {
		Repo       *user.Repo
		Commandbus commandbus.CommandBus
		Logger     logger.Logger
	}

	// RegistrationRequestBody user to register new users
	RegistrationRequestBody struct {
		Email string `json:"email"`
	}
)

// Register new user
func (c *CommandAPI) Register(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	opt := &RegistrationRequestBody{}
	if err := api.UnmarshalInto(body, opt); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	uID, err := uuid.NewV4()
	if err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	token, err := uuid.NewV4()
	if err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	tokenHash := hash(token.String())
	user := user.User{
		Metadata: user.Metadata{
			Email: opt.Email,
			ID:    uID.String(),
		},
		Spec: user.Spec{
			TokenHash: tokenHash,
		},
	}

	if err := c.Commandbus.ExecuteAndWait(ctx, "user.register", user); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}

	return api.NewAcceptedCommandResponseWithData("user", user.Metadata.ID, map[string]interface{}{
		"token": token.String(),
	})
}

func hash(s string) string {
	sh := sha256.Sum256([]byte(s))
	return string(sh[:])
}
