package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ngalaiko/miniboard/server/db"
	"github.com/stretchr/testify/assert"
)

func Test_Create(t *testing.T) {
	ctx := context.Background()

	sqlite := testDB(t)

	database := New(sqlite)

	testKey := &PublicKey{
		ID:        "test",
		DerBase64: "some data",
	}

	assert.NoError(t, database.Create(ctx, testKey))
}

func Test_Create_twice(t *testing.T) {
	ctx := context.Background()

	sqlite := testDB(t)

	database := New(sqlite)

	testKey := &PublicKey{
		ID:        "test",
		DerBase64: "some data",
	}

	assert.NoError(t, database.Create(ctx, testKey))
	assert.Error(t, database.Create(ctx, testKey))
}

func Test_Get(t *testing.T) {
	ctx := context.Background()

	sqlite := testDB(t)

	database := New(sqlite)

	testKey := &PublicKey{
		ID:        "test",
		DerBase64: "some data",
	}

	assert.NoError(t, database.Create(ctx, testKey))

	key, err := database.Get(ctx, testKey.ID)
	assert.NoError(t, err)

	assert.Equal(t, testKey.ID, key.ID)
	assert.Equal(t, testKey.DerBase64, key.DerBase64)
}

func Test_Get_not_exists(t *testing.T) {
	ctx := context.Background()

	sqlite := testDB(t)

	database := New(sqlite)

	testKey := &PublicKey{
		ID:        "test",
		DerBase64: "some data",
	}

	_, getErr := database.Get(ctx, testKey.ID)
	assert.Error(t, getErr)
}

func Test_Delete(t *testing.T) {
	ctx := context.Background()

	sqlite := testDB(t)

	database := New(sqlite)

	testKey := &PublicKey{
		ID:        "test",
		DerBase64: "some data",
	}

	assert.NoError(t, database.Create(ctx, testKey))

	_, getErr := database.Get(ctx, testKey.ID)
	assert.NoError(t, getErr)

	assert.NoError(t, database.Delete(ctx, testKey.ID))

	_, getAfterDeletion := database.Get(ctx, testKey.ID)
	assert.Error(t, getAfterDeletion)
}

func Test_Delete_not_existing(t *testing.T) {
	ctx := context.Background()

	sqlite := testDB(t)

	database := New(sqlite)

	testKey := &PublicKey{
		ID:        "test",
		DerBase64: "some data",
	}

	assert.NoError(t, database.Delete(ctx, testKey.ID))

	_, getAfterDeletion := database.Get(ctx, testKey.ID)
	assert.Error(t, getAfterDeletion)
}

func Test_List(t *testing.T) {
	ctx := context.Background()

	sqlite := testDB(t)

	database := New(sqlite)

	kk, err := database.List(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(kk))

	for i := 0; i < 10; i++ {
		testKey := &PublicKey{
			ID:        fmt.Sprintf("test-%d", i),
			DerBase64: "some data",
		}
		assert.NoError(t, database.Create(ctx, testKey))
	}

	kk, listError := database.List(ctx)
	assert.NoError(t, listError)
	assert.Equal(t, 10, len(kk))
}

func testDB(t *testing.T) *sql.DB {
	ctx := context.Background()

	tmpFile, err := ioutil.TempFile(os.TempDir(), "testdb-")
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, os.Remove(tmpFile.Name()))
	})

	db, err := db.New(ctx, &db.Config{
		Driver: "sqlite3",
		Addr:   tmpFile.Name(),
	}, &testLogger{})
	assert.NoError(t, err)
	return db
}

type testLogger struct{}

func (l *testLogger) Info(string, ...interface{}) {}
