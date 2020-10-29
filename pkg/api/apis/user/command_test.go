package user

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/peteqproj/peteq/pkg/utils/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommandAPI_Register(t *testing.T) {
	type fields struct {
		Repo        func() *user.Repo
		Commandbus  func() commandbus.CommandBus
		Logger      func() logger.Logger
		IDGenerator func() utils.IDGenerator
	}
	type args struct {
		ctx  context.Context
		body io.ReadCloser
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   api.CommandResponse
	}{
		{
			name: "Reject when request body does not match the validation schema",
			args: args{
				ctx:  context.Background(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{}),
			},
			want: api.NewRejectedCommandResponse(fmt.Errorf("Error: Email required | Error: Password required")),
		},
		{
			name: "Reject when request to execute command user.register failed",
			args: args{
				ctx: context.Background(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"email":    "peteq@peteq.io",
					"password": "123456",
				}),
			},
			fields: fields{
				IDGenerator: tests.NewIDBasicGenerator,
				Commandbus: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "user.register", mock.Anything).Return(fmt.Errorf("Error running command"))
					return cb
				},
				Logger: func() logger.Logger {
					l := &logger.MockLogger{}
					l.On("Info", "Failed to run user.register command", "error", "Error running command")
					return l
				},
			},
			want: api.NewRejectedCommandResponse(fmt.Errorf("Registration failed: Error running command")),
		},
		{
			name: "Reject when trying to register already exist email",
			args: args{
				ctx: context.Background(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"email":    "peteq@peteq.io",
					"password": "123456",
				}),
			},
			fields: fields{
				IDGenerator: tests.NewIDBasicGenerator,
				Logger: func() logger.Logger {
					l := &logger.MockLogger{}
					l.On("Info", "Failed to run user.register command", "error", "email exists")
					return l
				},
				Commandbus: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "user.register", mock.Anything).Return(fmt.Errorf("email exists"))
					return cb
				},
			},
			want: api.NewRejectedCommandResponse(fmt.Errorf("Registration failed: email exists")),
		},
		{
			name: "Register user",
			args: args{
				ctx: context.Background(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"email":    "peteq@peteq.io",
					"password": "123456",
				}),
			},
			fields: fields{
				IDGenerator: tests.NewIDBasicGenerator,
				Logger: func() logger.Logger {
					l := &logger.MockLogger{}
					return l
				},
				Commandbus: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "user.register", mock.Anything).Return(nil)
					return cb
				},
			},
			want: api.NewAcceptedCommandResponse("user", tests.GeneratedV4ID),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var repo *user.Repo
			var cm commandbus.CommandBus
			var logger logger.Logger
			var idGenerator utils.IDGenerator
			if tt.fields.Repo != nil {
				repo = tt.fields.Repo()
			}
			if tt.fields.Commandbus != nil {
				cm = tt.fields.Commandbus()
			}
			if tt.fields.Logger != nil {
				logger = tt.fields.Logger()
			}
			if tt.fields.IDGenerator != nil {
				idGenerator = tt.fields.IDGenerator()
			}
			ca := &CommandAPI{
				Repo:        repo,
				Commandbus:  cm,
				Logger:      logger,
				IDGenerator: idGenerator,
			}
			res := ca.Register(tt.args.ctx, tt.args.body)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestCommandAPI_Login(t *testing.T) {
	type fields struct {
		Repo        func() *user.Repo
		Commandbus  func() commandbus.CommandBus
		Logger      func() logger.Logger
		IDGenerator func() utils.IDGenerator
	}
	type args struct {
		ctx  context.Context
		body io.ReadCloser
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   api.CommandResponse
	}{
		{
			name: "Reject when request body does not match the validation schema",
			args: args{
				ctx:  context.Background(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{}),
			},
			want: api.NewRejectedCommandResponse(fmt.Errorf("Error: Email required | Error: Password required")),
		},
		{
			name: "Reject when request to execute command user.login failed",
			args: args{
				ctx: context.Background(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"email":    "peteq@peteq.io",
					"password": "123456",
				}),
			},
			fields: fields{
				IDGenerator: tests.NewIDBasicGenerator,
				Commandbus: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "user.login", mock.Anything).Return(fmt.Errorf("Error running command"))
					return cb
				},
				Logger: func() logger.Logger {
					l := &logger.MockLogger{}
					return l
				},
			},
			want: api.NewRejectedCommandResponse(fmt.Errorf("Login failed: Error running command")),
		},
		{
			name: "Reject when trying to login with invalid credentials",
			args: args{
				ctx: context.Background(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"email":    "peteq@peteq.io",
					"password": "123456",
				}),
			},
			fields: fields{
				IDGenerator: tests.NewIDBasicGenerator,
				Logger: func() logger.Logger {
					l := &logger.MockLogger{}
					return l
				},
				Commandbus: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "user.login", mock.Anything).Return(fmt.Errorf("invalid credentials"))
					return cb
				},
			},
			want: api.NewRejectedCommandResponse(fmt.Errorf("Login failed: invalid credentials")),
		},
		{
			name: "Login",
			args: args{
				ctx: context.Background(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"email":    "peteq@peteq.io",
					"password": "123456",
				}),
			},
			fields: fields{
				IDGenerator: tests.NewIDBasicGenerator,
				Logger: func() logger.Logger {
					l := &logger.MockLogger{}
					return l
				},
				Commandbus: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "user.login", mock.Anything).Return(nil)
					return cb
				},
			},
			want: api.NewAcceptedCommandResponseWithData("user", "", map[string]string{
				"token": tests.GeneratedV4ID,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var repo *user.Repo
			var cm commandbus.CommandBus
			var logger logger.Logger
			var idGenerator utils.IDGenerator
			if tt.fields.Repo != nil {
				repo = tt.fields.Repo()
			}
			if tt.fields.Commandbus != nil {
				cm = tt.fields.Commandbus()
			}
			if tt.fields.Logger != nil {
				logger = tt.fields.Logger()
			}
			if tt.fields.IDGenerator != nil {
				idGenerator = tt.fields.IDGenerator()
			}
			c := &CommandAPI{
				Repo:        repo,
				Commandbus:  cm,
				Logger:      logger,
				IDGenerator: idGenerator,
			}
			res := c.Login(tt.args.ctx, tt.args.body)
			assert.Equal(t, tt.want, res)
		})
	}
}
