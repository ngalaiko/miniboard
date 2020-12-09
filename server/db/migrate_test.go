package db

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Migrate(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "testdb-")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	db, err := New(&Config{
		Driver: "sqlite3",
		Addr:   tmpFile.Name(),
	}, &testLogger{})
	assert.NoError(t, err)

	assert.NoError(t, Migrate(context.Background(), db, &testLogger{}))
}

func Test_Migrate_twice(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "testdb-")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	db, err := New(&Config{
		Driver: "sqlite3",
		Addr:   tmpFile.Name(),
	}, &testLogger{})
	assert.NoError(t, err)

	assert.NoError(t, Migrate(context.Background(), db, &testLogger{}))
	assert.NoError(t, Migrate(context.Background(), db, &testLogger{}))
}

type testLogger struct{}

func (l *testLogger) Info(string, ...interface{}) {}
