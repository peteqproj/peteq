package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/peteqproj/peteq/domain/user"
	"github.com/peteqproj/peteq/pkg/tenant"
	"github.com/peteqproj/peteq/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	// GeneratedV4ID generated id to be used as test id
	GeneratedV4ID = "776a866c-80c3-476e-b1d1-680c2296c43c"
)

// AuthenticatedContext creates authenticated context with user
func AuthenticatedContext() context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, tenant.User, user.User{
		Metadata: user.Metadata{
			ID: "user-id",
		},
		Spec: user.Spec{
			Email: "some@email.com",
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

// NewIDBasicGenerator common id generator
func NewIDBasicGenerator() utils.IDGenerator {
	i := &utils.MockIDGenerator{}
	i.On("GenerateV4").Return(GeneratedV4ID, nil)
	return i
}

// BuildDBConnectionOrDie created mocked database
func BuildDBConnectionOrDie() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	utils.DieOnError(err, "Failed to create sqlmock")
	pg := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gdb, err := gorm.Open(pg, &gorm.Config{})
	utils.DieOnError(err, "Failed to connect to mocked database")
	return gdb, mock
}
