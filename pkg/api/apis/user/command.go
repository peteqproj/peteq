package user

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"time"

	"github.com/gofrs/uuid"
	listCommand "github.com/peteqproj/peteq/domain/list/command"
	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/domain/user/command"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/tenant"
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// LoginRequestBody user to register new users
	LoginRequestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	// TODO: validate request
	if err := c.Commandbus.ExecuteAndWait(ctx, "user.register", command.RegisterCommandOptions{
		Email:        opt.Email,
		UserID:       uID.String(),
		PasswordHash: hash(opt.Password),
	}); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}

	basicLists := []string{"Upcoming", "Today", "Done"}
	ectx := tenant.ContextWithUser(ctx, user.User{
		Metadata: user.Metadata{
			Email: opt.Email,
			ID:    uID.String(),
		},
	})
	for i, l := range basicLists {
		time.Sleep(time.Second * 2)
		id, err := uuid.NewV4()
		if err != nil {
			return api.NewRejectedCommandResponse(err.Error())
		}
		if err := c.Commandbus.ExecuteAndWait(ectx, "list.create", listCommand.CreateCommandOptions{
			Name:  l,
			ID:    id.String(),
			Index: i,
		}); err != nil {
			return api.NewRejectedCommandResponse(err.Error())
		}
	}

	return api.NewAcceptedCommandResponse("user", uID.String())
}

// Login validates user exists and returns api token
func (c *CommandAPI) Login(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	opt := &LoginRequestBody{}
	if err := api.UnmarshalInto(body, opt); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	users, err := c.Repo.List(user.ListOptions{})
	if err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	validUserIndex := -1
	for i, u := range users {
		if u.Metadata.Email != opt.Email {
			continue
		}
		if hash(opt.Password) != u.Spec.PasswordHash {
			continue
		}
		validUserIndex = i
	}

	if validUserIndex == -1 {
		return api.NewRejectedCommandResponse("Invalid credentials")
	}

	token, err := uuid.NewV4()
	if err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}
	tokenHash := hash(token.String())
	if err := c.Commandbus.ExecuteAndWait(ctx, "user.login", command.LoginCommandOptions{
		HashedToken: tokenHash,
		UserID:      users[validUserIndex].Metadata.ID,
	}); err != nil {
		return api.NewRejectedCommandResponse(err.Error())
	}

	return api.NewAcceptedCommandResponseWithData("user", users[validUserIndex].Metadata.ID, map[string]string{
		"token": token.String(),
	})
}

func hash(s string) string {
	sh := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x\n", sh)
}
