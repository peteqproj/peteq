package list

import (
	"context"
	"fmt"
	"io"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/peteqproj/peteq/domain/list"
	"github.com/peteqproj/peteq/pkg/api"
	commandbus "github.com/peteqproj/peteq/pkg/command/bus"
	"github.com/peteqproj/peteq/pkg/logger"
	"github.com/peteqproj/peteq/pkg/utils/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommandAPI_MoveTasks(t *testing.T) {
	type fields struct {
		Repo       func() *list.Repo
		Commandbus func() commandbus.CommandBus
		Logger     func() logger.Logger
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
			name: "Reject when validation rejected",
			args: args{
				ctx:  tests.AuthenticatedContext(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{}),
			},
			want: api.NewRejectedCommandResponse(
				fmt.Errorf("Error: TaskIDs required"),
			),
		},
		{
			name: "Reject when source list not found",
			args: args{
				ctx: tests.AuthenticatedContext(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"source":      "not-found",
					"destination": "not-found",
					"tasks":       []string{"1"},
				}),
			},
			fields: fields{
				Repo: func() *list.Repo {
					db, mock, _ := sqlmock.New()
					q := ".*"
					mock.
						ExpectQuery(q).
						WillReturnRows(sqlmock.NewRows([]string{
							"userid",
							"listid",
							"info",
						}))
					return &list.Repo{
						DB: db,
					}
				},
			},
			want: api.NewRejectedCommandResponse(
				fmt.Errorf("Source list: List not found"),
			),
		},
		{
			name: "Reject when destination list not found",
			args: args{
				ctx: tests.AuthenticatedContext(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"source":      "found",
					"destination": "not-found",
					"tasks":       []string{"1"},
				}),
			},
			fields: fields{
				Repo: func() *list.Repo {
					l := string(tests.MustMarshal(list.List{
						Metadata: list.Metadata{
							ID:   "found",
							Name: "list",
						},
					}))
					db, mock, _ := sqlmock.New()
					q := ".*"
					mock.
						ExpectQuery(q).
						WillReturnRows(sqlmock.NewRows([]string{
							"userid",
							"listid",
							"info",
						}).AddRow("userid", "listid", l))
					mock.
						ExpectQuery(q).
						WillReturnRows(sqlmock.NewRows([]string{
							"userid",
							"listid",
							"info",
						}))

					return &list.Repo{
						DB: db,
					}
				},
			},
			want: api.NewRejectedCommandResponse(
				fmt.Errorf("Destination list: List not found"),
			),
		},
		{
			name: "Reject when call the commandbus failed on command list.move-task",
			args: args{
				ctx: tests.AuthenticatedContext(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"source":      "source",
					"destination": "destination",
					"tasks":       []string{"1"},
				}),
			},
			fields: fields{
				Repo: func() *list.Repo {
					db, mock, _ := sqlmock.New()
					q := ".*"
					mock.
						ExpectQuery(q).
						WillReturnRows(sqlmock.NewRows([]string{
							"userid",
							"listid",
							"info",
						}).AddRow("userid", "listid", string(tests.MustMarshal(list.List{
							Metadata: list.Metadata{
								ID:   "source",
								Name: "source",
							},
						}))))
					mock.
						ExpectQuery(q).
						WillReturnRows(sqlmock.NewRows([]string{
							"userid",
							"listid",
							"info",
						}).AddRow("userid", "listid", string(tests.MustMarshal(list.List{
							Metadata: list.Metadata{
								ID:   "destination",
								Name: "destination",
							},
						}))))

					return &list.Repo{
						DB: db,
					}
				},
				Commandbus: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "list.move-task", mock.Anything).Return(fmt.Errorf("Failed to execute command list.move-task"))
					return cb
				},
				Logger: func() logger.Logger {
					l := &logger.MockLogger{}
					l.On("Info", "Moving task", "source", "source", "destination", "destination", "task", "1")
					l.On("Info", "Failed to execute command list.move-task", "error", "Failed to execute command list.move-task")
					return l
				},
			},
			want: api.NewRejectedCommandResponse(fmt.Errorf("Failed to move task 1")),
		},
		{
			name: "Reject when call the commandbus failed on command task.complete",
			args: args{
				ctx: tests.AuthenticatedContext(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"source":      "Upcoming",
					"destination": "Done",
					"tasks":       []string{"1"},
				}),
			},
			fields: fields{
				Repo: func() *list.Repo {
					db, mock, _ := sqlmock.New()
					q := ".*"
					mock.
						ExpectQuery(q).
						WillReturnRows(sqlmock.NewRows([]string{
							"userid",
							"listid",
							"info",
						}).AddRow("userid", "listid", string(tests.MustMarshal(list.List{
							Metadata: list.Metadata{
								ID:   "Upcoming",
								Name: "Upcoming",
							},
						}))))
					mock.
						ExpectQuery(q).
						WillReturnRows(sqlmock.NewRows([]string{
							"userid",
							"listid",
							"info",
						}).AddRow("userid", "listid", string(tests.MustMarshal(list.List{
							Metadata: list.Metadata{
								ID:   "Done",
								Name: "Done",
							},
						}))))

					return &list.Repo{
						DB: db,
					}
				},
				Commandbus: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "list.move-task", mock.Anything).Return(nil)
					cb.On("Execute", mock.Anything, "task.complete", mock.Anything).Return(fmt.Errorf("Failed to execute command"))
					return cb
				},
				Logger: func() logger.Logger {
					l := &logger.MockLogger{}
					l.On("Info", "Moving task", "source", "Upcoming", "destination", "Done", "task", "1")
					l.On("Info", "Failed to execute command task.complete", "error", "Failed to execute command")
					l.On("Info", "Completing task", "name", "1")
					return l
				},
			},
			want: api.NewRejectedCommandResponse(fmt.Errorf("Failed to complete task 1")),
		},
		{
			name: "Reject when call the commandbus failed on command task.reopen",
			args: args{
				ctx: tests.AuthenticatedContext(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"source":      "Done",
					"destination": "Upcoming",
					"tasks":       []string{"1"},
				}),
			},
			fields: fields{
				Repo: func() *list.Repo {
					db, mock, _ := sqlmock.New()
					q := ".*"
					mock.
						ExpectQuery(q).
						WillReturnRows(sqlmock.NewRows([]string{
							"userid",
							"listid",
							"info",
						}).AddRow("userid", "listid", string(tests.MustMarshal(list.List{
							Metadata: list.Metadata{
								ID:   "Done",
								Name: "Done",
							},
						}))))
					mock.
						ExpectQuery(q).
						WillReturnRows(sqlmock.NewRows([]string{
							"userid",
							"listid",
							"info",
						}).AddRow("userid", "listid", string(tests.MustMarshal(list.List{
							Metadata: list.Metadata{
								ID:   "Upcoming",
								Name: "Upcoming",
							},
						}))))

					return &list.Repo{
						DB: db,
					}
				},
				Commandbus: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "list.move-task", mock.Anything).Return(nil)
					cb.On("Execute", mock.Anything, "task.reopen", mock.Anything).Return(fmt.Errorf("Failed to execute command"))
					return cb
				},
				Logger: func() logger.Logger {
					l := &logger.MockLogger{}
					l.On("Info", "Moving task", "source", "Done", "destination", "Upcoming", "task", "1")
					l.On("Info", "Failed to execute command task.reopen", "error", "Failed to execute command")
					l.On("Info", "Reopenning task", "name", "1")
					return l
				},
			},
			want: api.NewRejectedCommandResponse(fmt.Errorf("Failed to reopen task 1")),
		},
		{
			name: "Accepct and move multiple tasks to destination list",
			args: args{
				ctx: tests.AuthenticatedContext(),
				body: tests.JSONStringToReadCloser(map[string]interface{}{
					"source":      "Upcoming",
					"destination": "Today",
					"tasks":       []string{"1", "2"},
				}),
			},
			fields: fields{
				Repo: func() *list.Repo {
					db, mock, _ := sqlmock.New()
					q := ".*"
					mock.
						ExpectQuery(q).
						WillReturnRows(sqlmock.NewRows([]string{
							"userid",
							"listid",
							"info",
						}).AddRow("userid", "listid", string(tests.MustMarshal(list.List{
							Metadata: list.Metadata{
								ID:   "Upcoming",
								Name: "Upcoming",
							},
						}))))
					mock.
						ExpectQuery(q).
						WillReturnRows(sqlmock.NewRows([]string{
							"userid",
							"listid",
							"info",
						}).AddRow("userid", "listid", string(tests.MustMarshal(list.List{
							Metadata: list.Metadata{
								ID:   "Today",
								Name: "Today",
							},
						}))))

					return &list.Repo{
						DB: db,
					}
				},
				Commandbus: func() commandbus.CommandBus {
					cb := &commandbus.MockCommandBus{}
					cb.On("Execute", mock.Anything, "list.move-task", mock.Anything).Return(nil).Times(2)
					return cb
				},
				Logger: func() logger.Logger {
					l := &logger.MockLogger{}
					l.On("Info", "Moving task", "source", "Upcoming", "destination", "Today", "task", "1")
					l.On("Info", "Moving task", "source", "Upcoming", "destination", "Today", "task", "2")
					return l
				},
			},
			want: api.NewAcceptedCommandResponse("list", "Upcoming"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var repo *list.Repo
			var cm commandbus.CommandBus
			var logger logger.Logger
			if tt.fields.Repo != nil {
				repo = tt.fields.Repo()
			}
			if tt.fields.Commandbus != nil {
				cm = tt.fields.Commandbus()
			}
			if tt.fields.Logger != nil {
				logger = tt.fields.Logger()
			}
			ca := &CommandAPI{
				Repo:       repo,
				Commandbus: cm,
				Logger:     logger,
			}
			resp := ca.MoveTasks(tt.args.ctx, tt.args.body)
			assert.Equal(t, tt.want, resp)
		})
	}
}
