package articles

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/ngalaiko/miniboard/server/db"
	"github.com/stretchr/testify/assert"
)

func Test_db_Create(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))
	testArticle := article()
	assert.NoError(t, database.Create(ctx, testArticle))
}

func Test_db_Create_twice(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	testArticle := article()

	assert.NoError(t, database.Create(ctx, testArticle))
	assert.Error(t, database.Create(ctx, testArticle))
}

func Test_db_Get(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	testArticle := article()

	assert.NoError(t, database.Create(ctx, testArticle))

	article, err := database.Get(ctx, testArticle.Id)
	assert.NoError(t, err)
	assert.Equal(t, testArticle, article)
}

func Test_db_Get_not_exists(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	testArticle := article()

	article, err := database.Get(ctx, testArticle.Id)
	assert.Error(t, err)
	assert.Nil(t, article)
}

func Test_db_GetByUserIDUrl(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	testArticle := article()

	assert.NoError(t, database.Create(ctx, testArticle))

	article, err := database.GetByUserIDUrl(ctx, testArticle.UserId, testArticle.Url)
	assert.NoError(t, err)

	assert.Equal(t, testArticle, article)
}

func Test_db_GetByUserIDUrl_not_exists(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	testArticle := article()

	article, err := database.GetByUserIDUrl(ctx, testArticle.UserId, testArticle.Url)
	assert.Error(t, err)
	assert.Nil(t, article)
}

func Test_db_Delete(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	testArticle := article()

	assert.NoError(t, database.Create(ctx, testArticle))

	article, err := database.Get(ctx, testArticle.Id)
	assert.NoError(t, err)
	assert.Equal(t, testArticle, article)

	assert.NoError(t, database.Delete(ctx, testArticle.Id))

	article, err = database.Get(ctx, testArticle.Id)
	assert.Error(t, err)
	assert.Nil(t, article)
}

func Test_db_Delete_not_existing(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	testArticle := article()

	assert.NoError(t, database.Delete(ctx, testArticle.Id))
}

func Test_db_Update_is_read(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	testArticle := article()

	assert.NoError(t, database.Create(ctx, testArticle))

	article, err := database.Get(ctx, testArticle.Id)
	assert.NoError(t, err)
	assert.Equal(t, testArticle, article)

	testArticle.IsRead = false
	assert.NoError(t, database.Update(ctx, testArticle))

	article, err = database.Get(ctx, testArticle.Id)
	assert.NoError(t, err)
	assert.Equal(t, testArticle, article)
}

func Test_db_List_all(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	saved := []*Article{}
	for i := 0; i < 10; i++ {
		a := article()
		a.Id += fmt.Sprint(i)
		a.ContentSha256 += fmt.Sprint(i)
		saved = append(saved, a)
		assert.NoError(t, database.Create(ctx, a))

		a.Content = nil
	}

	aa, err := database.List(ctx, &ListArticlesRequest{
		UserId:   saved[0].UserId,
		PageSize: 100,
	})
	assert.NoError(t, err)
	assert.Equal(t, 10, len(aa))
	assert.Equal(t, saved[0], aa[0])
	assert.Equal(t, saved[9], aa[9])
}

func Test_db_List_with_limit(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	saved := []*Article{}
	for i := 0; i < 10; i++ {
		a := article()
		a.Id += fmt.Sprint(i)
		a.ContentSha256 += fmt.Sprint(i)
		saved = append(saved, a)
		assert.NoError(t, database.Create(ctx, a))

		a.Content = nil
	}

	aa, err := database.List(ctx, &ListArticlesRequest{
		UserId:   saved[0].UserId,
		PageSize: 5,
	})
	assert.NoError(t, err)
	assert.Equal(t, 5, len(aa))
	assert.Equal(t, saved[0], aa[0])
	assert.Equal(t, saved[4], aa[4])
}

func Test_db_List_with_from(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	saved := []*Article{}
	for i := 0; i < 10; i++ {
		a := article()
		a.Id += fmt.Sprint(i)
		a.ContentSha256 += fmt.Sprint(i)
		saved = append(saved, a)
		assert.NoError(t, database.Create(ctx, a))

		a.Content = nil
	}

	aa, err := database.List(ctx, &ListArticlesRequest{
		UserId:    saved[0].UserId,
		PageToken: base64.StdEncoding.EncodeToString([]byte(saved[4].Id)),
		PageSize:  100,
	})
	assert.NoError(t, err)
	assert.Equal(t, 6, len(aa))
	assert.Equal(t, saved[4], aa[0])
	assert.Equal(t, saved[8], aa[4])
}

func Test_db_List_with_is_read(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	saved := []*Article{}
	for i := 0; i < 10; i++ {
		a := article()
		a.Id += fmt.Sprint(i)
		a.IsRead = i%2 == 0
		a.ContentSha256 += fmt.Sprint(i)
		saved = append(saved, a)
		assert.NoError(t, database.Create(ctx, a))

		a.Content = nil
	}

	aa, err := database.List(ctx, &ListArticlesRequest{
		UserId:   saved[0].UserId,
		PageSize: 100,
		IsReadEq: &wrappers.BoolValue{
			Value: true,
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 5, len(aa))
	for _, a := range aa {
		assert.Equal(t, true, a.IsRead)
	}
}

func Test_db_List_with_feed_id(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	saved := []*Article{}
	for i := 0; i < 10; i++ {
		a := article()
		a.Id += fmt.Sprint(i)
		a.ContentSha256 += fmt.Sprint(i)
		if i%2 == 0 {
			a.FeedId = "test"
		}
		saved = append(saved, a)
		assert.NoError(t, database.Create(ctx, a))

		a.Content = nil
	}

	aa, err := database.List(ctx, &ListArticlesRequest{
		UserId:   saved[0].UserId,
		PageSize: 100,
		FeedIdEq: &wrappers.StringValue{
			Value: "test",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 5, len(aa))
	for _, a := range aa {
		assert.Equal(t, "test", a.FeedId)
	}
}

func Test_db_List_with_title_contains(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	saved := []*Article{}
	for i := 0; i < 10; i++ {
		a := article()
		a.Id += fmt.Sprint(i)
		a.ContentSha256 += fmt.Sprint(i)
		if i%2 == 0 {
			a.Title = "it contains the query"
		}
		saved = append(saved, a)
		assert.NoError(t, database.Create(ctx, a))

		a.Content = nil
	}

	aa, err := database.List(ctx, &ListArticlesRequest{
		UserId:   saved[0].UserId,
		PageSize: 100,
		TitleContains: &wrappers.StringValue{
			Value: "contains the",
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 5, len(aa))
	for _, a := range aa {
		assert.Equal(t, "it contains the query", a.Title)
	}
}

func Test_db_List_with_title_contains_and_is_read(t *testing.T) {
	ctx := context.Background()
	database := newDB(testDB(t))

	saved := []*Article{}
	for i := 0; i < 10; i++ {
		a := article()
		a.Id += fmt.Sprint(i)
		a.ContentSha256 += fmt.Sprint(i)
		if i%2 == 0 {
			a.Title = "it contains the query"
			a.IsRead = true
		}
		saved = append(saved, a)
		assert.NoError(t, database.Create(ctx, a))

		a.Content = nil
	}

	aa, err := database.List(ctx, &ListArticlesRequest{
		UserId:   saved[0].UserId,
		PageSize: 100,
		TitleContains: &wrappers.StringValue{
			Value: "contains the",
		},
		IsReadEq: &wrappers.BoolValue{
			Value: true,
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 5, len(aa))
	for _, a := range aa {
		assert.Equal(t, "it contains the query", a.Title)
		assert.Equal(t, true, a.IsRead)
	}
}

func testDB(t *testing.T) *sql.DB {
	ctx := context.Background()

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

func article() *Article {
	return &Article{
		Id:            "id",
		UserId:        "user_id",
		Url:           "url",
		Title:         "title",
		CreateTime:    ptypes.TimestampNow(),
		Content:       []byte("content"),
		ContentSha256: "shasum",
		IsRead:        true,
	}
}
