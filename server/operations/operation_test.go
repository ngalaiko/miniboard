package operations

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/ngalaiko/miniboard/server/actor"
	"github.com/ngalaiko/miniboard/server/db"
	longrunning "github.com/ngalaiko/miniboard/server/genproto/google/longrunning"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func Test_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(testContext())
	defer cancel()

	o := New(testDB(ctx, t))

	runningOperations := make(chan *longrunning.Operation)
	f := func(ctx context.Context, operation *longrunning.Operation) (*any.Any, error) {
		runningOperations <- operation
		return &any.Any{}, nil
	}

	operation, err := o.CreateOperation(ctx, &any.Any{}, f)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(operation, <-runningOperations))
}

func Test_Get_done(t *testing.T) {
	ctx, cancel := context.WithCancel(testContext())
	defer cancel()

	o := New(testDB(ctx, t))

	f := func(ctx context.Context, operation *longrunning.Operation) (*any.Any, error) {
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

func Test_Get_error(t *testing.T) {
	ctx, cancel := context.WithCancel(testContext())
	defer cancel()

	o := New(testDB(ctx, t))

	f := func(ctx context.Context, operation *longrunning.Operation) (*any.Any, error) {
		return nil, fmt.Errorf("test error")
	}

	operation, err := o.CreateOperation(ctx, &any.Any{}, f)
	assert.NoError(t, err)
	assert.Eventually(t, func() bool {
		doneOperation, err := o.GetOperation(ctx, &longrunning.GetOperationRequest{
			Name: operation.Name,
		})
		assert.NoError(t, err)

		return doneOperation.GetError() != nil
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
