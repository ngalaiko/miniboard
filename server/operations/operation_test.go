package operations

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/ngalaiko/miniboard/server/actor"
	"github.com/ngalaiko/miniboard/server/db"
	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/genproto/googleapis/rpc/status"
)

func Test_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(testContext())
	defer cancel()

	o := New(testDB(ctx, t))

	runningOperations := make(chan *longrunning.Operation)
	f := func(ctx context.Context, operation *longrunning.Operation) (*any.Any, *status.Status) {
		runningOperations <- operation
		return &any.Any{}, nil
	}

	operation, err := o.CreateOperation(ctx, &any.Any{}, f)
	assert.NoError(t, err)
	assert.Equal(t, operation, <-runningOperations)
}

func Test_Get_done(t *testing.T) {
	ctx, cancel := context.WithCancel(testContext())
	defer cancel()

	o := New(testDB(ctx, t))

	f := func(ctx context.Context, operation *longrunning.Operation) (*any.Any, *status.Status) {
		return &any.Any{}, nil
	}

	operation, err := o.CreateOperation(ctx, &any.Any{}, f)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		doneOperation, err := o.GetOperation(ctx, &longrunning.GetOperationRequest{
			Name: operation.Name,
		})
		assert.NoError(t, err)

		return doneOperation.Done
	}, time.Second, 10*time.Millisecond)
}

func testContext() context.Context {
	return actor.NewContext(context.Background(), "test")
}

func testDB(ctx context.Context, t *testing.T) *sql.DB {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "testdb-")
	assert.NoError(t, err)

	go func() {
		<-ctx.Done()
		os.Remove(tmpFile.Name())
	}()

	sqlite, err := db.NewSQLite(tmpFile.Name())
	assert.NoError(t, err)
	assert.NoError(t, db.Migrate(ctx, sqlite))

	return sqlite
}
