package db

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/ngalaiko/miniboard/server/db"
	longrunning "github.com/ngalaiko/miniboard/server/genproto/google/longrunning"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/status"
)

func Test_Create(t *testing.T) {
	ctx := context.Background()
	database := New(testDB(t))
	assert.NoError(t, database.Create(ctx, "user", operation()))
}

func Test_Create_twice(t *testing.T) {
	ctx := context.Background()
	database := New(testDB(t))
	assert.NoError(t, database.Create(ctx, "user", operation()))
	assert.Error(t, database.Create(ctx, "user", operation()))
}

func Test_Get(t *testing.T) {
	ctx := context.Background()
	database := New(testDB(t))
	testOperation := operation()
	assert.NoError(t, database.Create(ctx, "user", testOperation))

	operation, err := database.Get(ctx, strings.Replace(testOperation.Name, "operations/", "", -1), "user")
	assert.NoError(t, err)
	assert.True(t, proto.Equal(testOperation, operation))
}

func Test_Get_not_exists(t *testing.T) {
	ctx := context.Background()
	database := New(testDB(t))

	testOperation := operation()
	operation, err := database.Get(ctx, strings.Replace(testOperation.Name, "operations/", "", -1), "user")
	assert.Nil(t, operation)
	assert.Equal(t, sql.ErrNoRows, err)
}

func Test_Update_result(t *testing.T) {
	ctx := context.Background()
	database := New(testDB(t))

	testOperation := operation()
	assert.NoError(t, database.Create(ctx, "user", testOperation))

	m, _ := ptypes.MarshalAny(&longrunning.Operation{})
	testOperation.Result = &longrunning.Operation_Response{
		Response: m,
	}

	assert.NoError(t, database.Update(ctx, "user", testOperation))

	operation, err := database.Get(ctx, strings.Replace(testOperation.Name, "operations/", "", -1), "user")
	assert.NoError(t, err)
	assert.True(t, proto.Equal(testOperation, operation))
}

func Test_Update_error(t *testing.T) {
	ctx := context.Background()
	database := New(testDB(t))

	testOperation := operation()
	assert.NoError(t, database.Create(ctx, "user", testOperation))

	testOperation.Result = &longrunning.Operation_Error{
		Error: &status.Status{
			Message: "test",
		},
	}

	assert.NoError(t, database.Update(ctx, "user", testOperation))

	operation, err := database.Get(ctx, strings.Replace(testOperation.Name, "operations/", "", -1), "user")
	assert.NoError(t, err)
	assert.True(t, proto.Equal(testOperation, operation))
}

func Test_Update_not_exists(t *testing.T) {
	ctx := context.Background()
	database := New(testDB(t))

	testOperation := operation()

	m, _ := ptypes.MarshalAny(&longrunning.Operation{})
	testOperation.Result = &longrunning.Operation_Response{
		Response: m,
	}

	assert.Equal(t, sql.ErrNoRows, database.Update(ctx, "user", testOperation))
}

func operation() *longrunning.Operation {
	m, _ := ptypes.MarshalAny(&longrunning.Operation{})
	return &longrunning.Operation{
		Name:     "operations/test",
		Done:     false,
		Metadata: m,
	}
}

func testDB(t *testing.T) *sql.DB {
	ctx := context.Background()

	tmpFile, err := ioutil.TempFile(os.TempDir(), "testdb-")
	assert.NoError(t, err)

	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	sqlite, err := db.New(ctx, &db.Config{
		Driver: "sqlite3",
		Addr:   tmpFile.Name(),
	}, &testLogger{})
	assert.NoError(t, err)

	return sqlite
}

type testLogger struct{}

func (l *testLogger) Info(string, ...interface{}) {}
