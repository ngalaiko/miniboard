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

	_, err = New(context.Background(), &Config{
		Driver: "sqlite3",
		Addr:   tmpFile.Name(),
	})
	assert.NoError(t, err)
}

func Test_Migrate_twice(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "testdb-")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = New(context.Background(), &Config{
		Driver: "sqlite3",
		Addr:   tmpFile.Name(),
	})
	assert.NoError(t, err)

	_, err = New(context.Background(), &Config{
		Driver: "sqlite3",
		Addr:   tmpFile.Name(),
	})
	assert.NoError(t, err)
}
