package db

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/ngalaiko/miniboard/server/actor"
	"github.com/ngalaiko/miniboard/server/db"
	"github.com/ngalaiko/miniboard/server/genproto/feeds/v1"
	"github.com/stretchr/testify/assert"
)

func Test_db_Create(t *testing.T) {
	ctx := testContext()
	database := New(testDB(t))
	testFeed := feed()
	assert.NoError(t, database.Create(ctx, testFeed))
}

func Test_db_Create_twice(t *testing.T) {
	ctx := testContext()
	database := New(testDB(t))
	testFeed := feed()
	assert.NoError(t, database.Create(ctx, testFeed))
	assert.Error(t, database.Create(ctx, testFeed))
}

func Test_db_Get(t *testing.T) {
	ctx := testContext()
	database := New(testDB(t))
	testFeed := feed()
	assert.NoError(t, database.Create(ctx, testFeed))

	feed, err := database.Get(ctx, testFeed.Id, testFeed.UserId)
	assert.NoError(t, err)
	assert.Equal(t, testFeed, feed)
}

func Test_db_Get_not_exists(t *testing.T) {
	ctx := testContext()
	database := New(testDB(t))
	testFeed := feed()

	feed, err := database.Get(ctx, testFeed.Id, testFeed.UserId)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.Nil(t, feed)
}

func Test_db_Update_timestamp(t *testing.T) {
	ctx := testContext()
	database := New(testDB(t))
	testFeed := feed()
	assert.NoError(t, database.Create(ctx, testFeed))

	testFeed.LastFetched = ptypes.TimestampNow()

	assert.NoError(t, database.Update(ctx, testFeed, testFeed.UserId))

	feed, err := database.Get(ctx, testFeed.Id, testFeed.UserId)
	assert.NoError(t, err)
	assert.Equal(t, testFeed, feed)
}

func Test_db_List_all(t *testing.T) {
	ctx := testContext()
	database := New(testDB(t))

	saved := []*feeds.Feed{}
	for i := 0; i < 10; i++ {
		f := feed()
		f.Id += fmt.Sprint(i)
		f.Url += fmt.Sprint(i)
		saved = append(saved, f)
		assert.NoError(t, database.Create(ctx, f))
	}

	ff, err := database.List(ctx, "user_id", &feeds.ListFeedsRequest{
		PageSize: 100,
	})
	assert.NoError(t, err)
	assert.Equal(t, 10, len(ff))
	assert.Equal(t, saved[0], ff[0])
	assert.Equal(t, saved[9], ff[9])
}

func Test_db_List_with_limit(t *testing.T) {
	ctx := testContext()
	database := New(testDB(t))

	saved := []*feeds.Feed{}
	for i := 0; i < 10; i++ {
		f := feed()
		f.Id += fmt.Sprint(i)
		f.Url += fmt.Sprint(i)
		saved = append(saved, f)
		assert.NoError(t, database.Create(ctx, f))
	}

	ff, err := database.List(ctx, "user_id", &feeds.ListFeedsRequest{
		PageSize: 5,
	})
	assert.NoError(t, err)
	assert.Equal(t, 5, len(ff))
	assert.Equal(t, saved[0], ff[0])
	assert.Equal(t, saved[4], ff[4])
}

func Test_db_List_with_from(t *testing.T) {
	ctx := testContext()
	database := New(testDB(t))

	saved := []*feeds.Feed{}
	for i := 0; i < 10; i++ {
		f := feed()
		f.Id += fmt.Sprint(i)
		f.Url += fmt.Sprint(i)
		saved = append(saved, f)
		assert.NoError(t, database.Create(ctx, f))
	}

	ff, err := database.List(ctx, "user_id", &feeds.ListFeedsRequest{
		PageToken: base64.StdEncoding.EncodeToString([]byte(saved[4].Id)),
		PageSize:  100,
	})
	assert.NoError(t, err)
	assert.Equal(t, 6, len(ff))
	assert.Equal(t, saved[4], ff[0])
	assert.Equal(t, saved[8], ff[4])
}

func Test_db_ListAll(t *testing.T) {
	ctx := testContext()
	database := New(testDB(t))

	saved := []*feeds.Feed{}
	for i := 0; i < 10; i++ {
		f := feed()
		f.Id += fmt.Sprint(i)
		f.UserId += fmt.Sprint(i)
		saved = append(saved, f)
		assert.NoError(t, database.Create(ctx, f))
	}

	ff, err := database.ListAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(ff))
	assert.Equal(t, saved[0], ff[0])
	assert.Equal(t, saved[9], ff[9])
}

func feed() *feeds.Feed {
	return &feeds.Feed{
		Id:          "id",
		UserId:      "user_id",
		LastFetched: ptypes.TimestampNow(),
		Url:         "url",
		Title:       " title",
	}
}

func testContext() context.Context {
	return actor.NewContext(context.Background(), "user_id")
}

func testDB(t *testing.T) *sql.DB {
	ctx := testContext()

	tmpFile, err := ioutil.TempFile(os.TempDir(), "testdb-")
	assert.NoError(t, err)

	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	sqlite, err := db.NewSQLite(tmpFile.Name())
	assert.NoError(t, err)
	assert.NoError(t, db.Migrate(ctx, sqlite))

	return sqlite
}
