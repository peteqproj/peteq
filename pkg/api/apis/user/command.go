package user

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/domain/user/command"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
)

type (
	// CommandAPI for users
	CommandAPI struct {
		Repo        *user.Repo
		Commandbus  commandbus.CommandBus
		Logger      logger.Logger
		IDGenerator utils.IDGenerator
	}

	// RegistrationRequestBody user to register new users
	RegistrationRequestBody struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	// LoginRequestBody user to register new users
	LoginRequestBody struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
)

// Register new user
func (c *CommandAPI) Register(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	opt := RegistrationRequestBody{}
	if err := api.UnmarshalInto(body, &opt); err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	uID, err := c.IDGenerator.GenerateV4()
	if err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	if err := c.Commandbus.Execute(ctx, "user.register", command.RegisterCommandOptions{
		Email:        opt.Email,
		UserID:       uID,
		PasswordHash: hash(opt.Password),
	}); err != nil {
		c.Logger.Info("Failed to run user.register command", "error", err.Error())
		return api.NewRejectedCommandResponse(fmt.Errorf("Registration failed: %v", err))
	}
	return api.NewAcceptedCommandResponse("user", uID)
}

// Login validates user exists and returns api token
func (c *CommandAPI) Login(ctx context.Context, body io.ReadCloser) api.CommandResponse {
	opt := &LoginRequestBody{}
	if err := api.UnmarshalInto(body, opt); err != nil {
		return api.NewRejectedCommandResponse(err)
	}

	token, err := c.IDGenerator.GenerateV4()
	if err != nil {
		return api.NewRejectedCommandResponse(err)
	}
	tokenHash := hash(token)
	if err := c.Commandbus.Execute(ctx, "user.login", command.LoginCommandOptions{
		HashedToken:    tokenHash,
		Email:          opt.Email,
		HashedPassword: hash(opt.Password),
	}); err != nil {
		return api.NewRejectedCommandResponse(fmt.Errorf("Login failed: %v", err))
	}

	return api.NewAcceptedCommandResponseWithData("user", "", map[string]string{
		"token": token,
	})
}

func hash(s string) string {
	sh := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x\n", sh)
}
