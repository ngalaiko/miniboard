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
)

func Test_Get_done(t *testing.T) {
	ctx, cancel := context.WithCancel(testContext())
	defer cancel()

	o := New(ctx, testDB(ctx, t))

	f := func(ctx context.Context, operation *longrunning.Operation, status chan<- *longrunning.Operation) error {
		operation.Done = true
		status <- operation
		return nil
	}

	operation, err := o.CreateOperation(ctx, &any.Any{}, f)
	assert.NoError(t, err)

	assert.Eventually(t, func() bool {
		doneOperation, err := o.GetOperation(ctx, &longrunning.GetOperationRequest{
			Name: operation.Name,
		})
		assert.NoError(t, err)

		return doneOperation.Done && doneOperation.GetError() == nil
	}, time.Second, 10*time.Millisecond)
}

func Test_Get_panic(t *testing.T) {
	ctx, cancel := context.WithCancel(testContext())
	defer cancel()

	o := New(ctx, testDB(ctx, t))

	f := func(ctx context.Context, operation *longrunning.Operation, status chan<- *longrunning.Operation) error {
		panic("test")
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

func Test_Get_error(t *testing.T) {
	ctx, cancel := context.WithCancel(testContext())
	defer cancel()

	o := New(ctx, testDB(ctx, t))

	f := func(ctx context.Context, operation *longrunning.Operation, status chan<- *longrunning.Operation) error {
		return fmt.Errorf("testError")
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

func Test_Get_no_updates(t *testing.T) {
	ctx, cancel := context.WithCancel(testContext())
	defer cancel()

	o := New(ctx, testDB(ctx, t))

	f := func(ctx context.Context, operation *longrunning.Operation, status chan<- *longrunning.Operation) error {
		return nil
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
