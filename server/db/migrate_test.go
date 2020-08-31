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

	db, err := NewSQLite(tmpFile.Name())
	assert.NoError(t, err)

	assert.NoError(t, Migrate(context.Background(), db))
}

func Test_Migrate_twice(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "testdb-")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	db, err := NewSQLite(tmpFile.Name())
	assert.NoError(t, err)

	assert.NoError(t, Migrate(context.Background(), db))
	assert.NoError(t, Migrate(context.Background(), db))
}
