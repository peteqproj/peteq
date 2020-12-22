package saga

import (
	"context"
	"testing"

	"github.com/peteqproj/peteq/domain/list"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/errors"
	"github.com/peteqproj/peteq/pkg/logger"
	tenant "github.com/peteqproj/peteq/pkg/tenant/test"
	"github.com/peteqproj/peteq/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errRunCommandBus = errors.New("Failed to run command bus")

func buildMockCommandBus() commandbus.CommandBus {
	cb := &commandbus.MockCommandBus{}
	cb.On("Execute", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	return cb
}

func buildMockLogger() logger.Logger {
	lgr := &logger.MockLogger{}
	lgr.On("Info", mock.Anything).Return(nil)
	return lgr
}

func buildMockIDGenerator() utils.IDGenerator {
	gen := &utils.MockIDGenerator{}
	gen.On("GenerateV4").Return("fake-id", nil)
	return gen
}

func buildMockListRepo() ListRepo {
	m := &MockListRepo{}
	m.On("GetListByName", mock.Anything, mock.Anything).Return(list.List{}, nil)
	return m
}

func Test_registrator_Run(t *testing.T) {
	type fields struct {
		CommandbusFn  func() commandbus.CommandBus
		LoggerFn      func() logger.Logger
		IDGeneratorFn func() utils.IDGenerator
		ListRepoFn    func() ListRepo
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		errorMessage string
	}{
		{
			name:         "Context is not authenticated -> return an error",
			errorMessage: "Authentication Error: user not found in context",
			fields: fields{
				CommandbusFn: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, mock.Anything, mock.Anything).Return(nil)
					return cb
				},
				LoggerFn:      buildMockLogger,
				IDGeneratorFn: buildMockIDGenerator,
				ListRepoFn:    buildMockListRepo,
			},
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name:         "Failed to create list -> return an error",
			errorMessage: "Failed to create list Upcoming: Failed to run command bus",
			fields: fields{
				CommandbusFn: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "list.create", mock.Anything).Return(errRunCommandBus)
					return cb
				},
				LoggerFn:      buildMockLogger,
				IDGeneratorFn: buildMockIDGenerator,
				ListRepoFn:    buildMockListRepo,
			},
			args: args{
				ctx: tenant.BuildAuthenticationContextWithUser(),
			},
		},
		{
			name:         "Failed to create trigger -> return an error",
			errorMessage: "Failed to create trigger Task Archiver: Failed to run command bus",
			fields: fields{
				CommandbusFn: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "list.create", mock.Anything).Return(nil)
					cb.On("Execute", mock.Anything, "trigger.create", mock.Anything).Return(errRunCommandBus)
					return cb
				},
				LoggerFn:      buildMockLogger,
				IDGeneratorFn: buildMockIDGenerator,
				ListRepoFn:    buildMockListRepo,
			},
			args: args{
				ctx: tenant.BuildAuthenticationContextWithUser(),
			},
		},
		{
			name:         "Failed to create automation -> return an error",
			errorMessage: "Failed to create automation Task Archiver: Failed to run command bus",
			fields: fields{
				CommandbusFn: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "list.create", mock.Anything).Return(nil)
					cb.On("Execute", mock.Anything, "trigger.create", mock.Anything).Return(nil)
					cb.On("Execute", mock.Anything, "automation.create", mock.Anything).Return(errRunCommandBus)
					return cb
				},
				LoggerFn:      buildMockLogger,
				IDGeneratorFn: buildMockIDGenerator,
				ListRepoFn:    buildMockListRepo,
			},
			args: args{
				ctx: tenant.BuildAuthenticationContextWithUser(),
			},
		},
		{
			name:         "Failed to create trigger-binding -> return an error",
			errorMessage: "Failed to automation-trigger-binding for Task Archiver trigger: Failed to run command bus",
			fields: fields{
				CommandbusFn: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "list.create", mock.Anything).Return(nil)
					cb.On("Execute", mock.Anything, "trigger.create", mock.Anything).Return(nil)
					cb.On("Execute", mock.Anything, "automation.create", mock.Anything).Return(nil)
					cb.On("Execute", mock.Anything, "automation.bindTrigger", mock.Anything).Return(errRunCommandBus)
					return cb
				},
				LoggerFn:      buildMockLogger,
				IDGeneratorFn: buildMockIDGenerator,
				ListRepoFn: func() ListRepo {
					m := &MockListRepo{}
					m.
						On("GetListByName", mock.Anything, mock.Anything).Return(list.List{}, errors.NewNotFoundError("List", "fake-id"))
					return m
				},
			},
			args: args{
				ctx: tenant.BuildAuthenticationContextWithUser(),
			},
		},
		{
			name: "In case a one of the list already exists do not attempt to create one",
			fields: fields{
				CommandbusFn:  buildMockCommandBus,
				LoggerFn:      buildMockLogger,
				IDGeneratorFn: buildMockIDGenerator,
				ListRepoFn:    buildMockListRepo,
			},
			args: args{
				ctx: tenant.BuildAuthenticationContextWithUser(),
			},
		},
		// {
		// 	name:    "Trigger already exists -> do not attempt to create one",
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Trigger-binding already exists -> do not attempt to create one",
		// 	wantErr: false,
		// },
		// {
		// 	name:    "Automation already exists -> do not attempt to create one",
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &registrator{
				Commandbus:  tt.fields.CommandbusFn(),
				Logger:      tt.fields.LoggerFn(),
				IDGenerator: tt.fields.IDGeneratorFn(),
				ListRepo:    tt.fields.ListRepoFn(),
			}
			err := a.Run(tt.args.ctx)
			if tt.errorMessage != "" {
				assert.EqualError(t, err, tt.errorMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
