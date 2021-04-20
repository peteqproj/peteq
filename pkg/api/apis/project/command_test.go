package project

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/peteqproj/peteq/domain/project"
	"github.com/peteqproj/peteq/domain/project/command"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/peteqproj/peteq/pkg/utils/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommandAPI_AddTasks(t *testing.T) {
	type fields struct {
		Repo        func() *project.Repo
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
			name: "Reject when request body is not matched to the schema",
			args: args{
				ctx:  tests.AuthenticatedContext(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{}),
			},
			want: api.NewRejectedCommandResponse(fmt.Errorf("Error: Project required | Error: TaskIDs required")),
		},
		{
			name: "Reject when execution of command failed",
			args: args{
				ctx: tests.AuthenticatedContext(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"project": "project",
					"tasks":   []string{"1", "2"},
				}),
			},
			fields: fields{
				Commandbus: func() commandbus.CommandBus {
					cm := &commandbus.MockCommandBus{}
					cm.On("Execute", mock.Anything, "project.add-task", mock.Anything).Return(fmt.Errorf("Failed to run command"))
					return cm
				},
				Logger: func() logger.Logger {
					l := &logger.MockLogger{}
					l.On("Info", "Failed to execute command project.add-task", "error", "Failed to run command")
					return l
				},
			},
			want: api.NewRejectedCommandResponse(fmt.Errorf("Failed to add task 1 to project project")),
		},
		{
			name: "Accept to add multiple tasks into project",
			args: args{
				ctx: tests.AuthenticatedContext(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"project": "project",
					"tasks":   []string{"1", "2"},
				}),
			},
			fields: fields{
				Commandbus: func() commandbus.CommandBus {
					cm := &commandbus.MockCommandBus{}
					cm.On("Execute", mock.Anything, "project.add-task", mock.Anything).Return(nil).Times(2)
					return cm
				},
			},
			want: api.NewAcceptedCommandResponse("project", "project"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var repo *project.Repo
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
			res := ca.AddTasks(tt.args.ctx, tt.args.body)
			assert.Equal(t, tt.want, res)
		})
	}
}

func TestCommandAPI_Create(t *testing.T) {
	type fields struct {
		Repo        func() *project.Repo
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
			name: "Reject when request body doest match to the schema",
			args: args{
				ctx:  tests.AuthenticatedContext(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{}),
			},
			want: api.NewRejectedCommandResponse(fmt.Errorf("Error: Name required")),
		},
		{
			name: "Reject when request to the command bus project.create failed",
			args: args{
				ctx: tests.AuthenticatedContext(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"name": "project",
				}),
			},
			fields: fields{
				Logger: func() logger.Logger {
					l := &logger.MockLogger{}
					l.On("Info", "Failed to execute project.create command", "error", "Failed to run command")
					return l
				},
				IDGenerator: func() utils.IDGenerator {
					i := &utils.MockIDGenerator{}
					i.On("GenerateV4").Return("project-id", nil)
					return i
				},
				Commandbus: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "project.create", mock.Anything).Return(fmt.Errorf("Failed to run command"))
					return cb
				},
			},
			want: api.NewRejectedCommandResponse(fmt.Errorf("Failed to create project")),
		},
		{
			name: "Accept command to create the project",
			args: args{
				ctx: tests.AuthenticatedContext(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"name": "project",
				}),
			},
			fields: fields{
				Logger: func() logger.Logger {
					l := &logger.MockLogger{}
					return l
				},
				Commandbus: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					opt := command.CreateProjectCommandOptions{
						ID:   "project-id",
						Name: "project",
					}
					cb.On("Execute", mock.Anything, "project.create", opt).Return(nil)
					return cb
				},
				IDGenerator: func() utils.IDGenerator {
					i := &utils.MockIDGenerator{}
					i.On("GenerateV4").Return("project-id", nil)
					return i
				},
			},
			want: api.NewAcceptedCommandResponse("project", "project-id"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var repo *project.Repo
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
			res := ca.Create(tt.args.ctx, tt.args.body)
			assert.Equal(t, tt.want, res)
		})
	}
}
