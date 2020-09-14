package articles

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/ngalaiko/miniboard/server/actor"
	"github.com/ngalaiko/miniboard/server/storage/resource"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func testArticle(replacement string) io.ReadCloser {
	file, err := os.Open("./testdata/test.html")
	if err != nil {
		return nil
	}

	dd, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}

	dd = bytes.Replace(dd, []byte("__RANDOM__"), []byte(replacement), 1)

	return ioutil.NopCloser(bytes.NewBuffer(dd))
}

type testClient struct{}

func (tc *testClient) Fetch(ctx context.Context, url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       testArticle(url),
	}, nil
}

func Test_service_Create(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userName := resource.NewName("users", "test")
	ctx = actor.NewContext(ctx, userName)

	service := NewService(testDB(t), &testClient{})

	testURL, _ := url.Parse("http://localhost")

	resp, err := service.CreateArticle(ctx, testArticle(testURL.String()), testURL, nil, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Id)
	assert.Equal(t, resp.Url, "http://localhost")
	assert.NotEmpty(t, resp.Content)
}

func Test_service_Create_twice_with_same_content(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userName := resource.NewName("users", "test")
	ctx = actor.NewContext(ctx, userName)

	service := NewService(testDB(t), &testClient{})

	testURL, _ := url.Parse("http://localhost")

	resp, err := service.CreateArticle(ctx, testArticle(testURL.String()), testURL, nil, "")
	assert.NoError(t, err)

	secondResponse, err := service.CreateArticle(ctx, testArticle(testURL.String()), testURL, nil, "")
	assert.NoError(t, err)
	assert.Equal(t, resp.CreateTime, secondResponse.CreateTime)
}

func Test_service_Create_twice_with_different_content(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userName := resource.NewName("users", "test")
	ctx = actor.NewContext(ctx, userName)

	service := NewService(testDB(t), &testClient{})

	testURL, _ := url.Parse("http://localhost")

	resp, err := service.CreateArticle(ctx, testArticle(testURL.String()), testURL, nil, "")
	assert.NoError(t, err)

	secondResponse, err := service.CreateArticle(ctx, testArticle("new content"), testURL, nil, "")
	assert.NoError(t, err)
	assert.Equal(t, resp.CreateTime, secondResponse.CreateTime)
	assert.Equal(t, resp.Id, secondResponse.Id)
	assert.Equal(t, resp.CreateTime, secondResponse.CreateTime)
	assert.NotEqual(t, resp.Content, secondResponse.Content)
	assert.NotEqual(t, resp.ContentSha256, secondResponse.ContentSha256)
}

func Test_service_Get_basic_view(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userName := resource.NewName("users", "test")
	ctx = actor.NewContext(ctx, userName)

	service := NewService(testDB(t), &testClient{})

	testURL, _ := url.Parse("http://localhost")

	resp, err := service.CreateArticle(ctx, testArticle(testURL.String()), testURL, nil, "")
	assert.NoError(t, err)

	article, err := service.GetArticle(ctx, &GetArticleRequest{
		Id: resp.Id,
	})
	assert.NoError(t, err)
	assert.Nil(t, article.Content)
}

func Test_service_Get_full_view(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userName := resource.NewName("users", "test")
	ctx = actor.NewContext(ctx, userName)

	service := NewService(testDB(t), &testClient{})

	testURL, _ := url.Parse("http://localhost")

	resp, err := service.CreateArticle(ctx, testArticle(testURL.String()), testURL, nil, "")
	assert.NoError(t, err)

	article, err := service.GetArticle(ctx, &GetArticleRequest{
		Id:   resp.Id,
		View: ArticleView_ARTICLE_VIEW_FULL,
	})
	assert.NoError(t, err)
	assert.Equal(t, resp, article)
}

func Test_service_Get_not_exists(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userName := resource.NewName("users", "test")
	ctx = actor.NewContext(ctx, userName)

	service := NewService(testDB(t), &testClient{})

	article, err := service.GetArticle(ctx, &GetArticleRequest{
		Id: "404",
	})
	assert.Nil(t, article)

	status, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, status.Code())
}

func Test_service_Delete(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userName := resource.NewName("users", "test")
	ctx = actor.NewContext(ctx, userName)

	service := NewService(testDB(t), &testClient{})

	testURL, _ := url.Parse("http://localhost")

	resp, err := service.CreateArticle(ctx, testArticle(testURL.String()), testURL, nil, "")
	assert.NoError(t, err)

	_, deleteErr := service.DeleteArticle(ctx, &DeleteArticleRequest{
		Id: resp.Id,
	})
	assert.NoError(t, deleteErr)

	_, getErr := service.GetArticle(ctx, &GetArticleRequest{
		Id: resp.Id,
	})

	status, ok := status.FromError(getErr)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, status.Code())
}

func Test_service_List_all(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userName := resource.NewName("users", "test")
	ctx = actor.NewContext(ctx, userName)

	service := NewService(testDB(t), &testClient{})

	for i := 0; i < 50; i++ {
		testURL, _ := url.Parse(fmt.Sprintf("http://localhost-%d", i))
		_, err := service.CreateArticle(ctx, testArticle(testURL.String()), testURL, nil, "")
		assert.NoError(t, err)
	}

	response, err := service.ListArticles(ctx, &ListArticlesRequest{
		PageSize: 100,
	})
	assert.NoError(t, err)
	assert.Equal(t, 50, len(response.Articles))
	assert.Empty(t, response.NextPageToken)
}

func Test_service_List_pagination(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userName := resource.NewName("users", "test")
	ctx = actor.NewContext(ctx, userName)

	service := NewService(testDB(t), &testClient{})

	for i := 0; i < 50; i++ {
		testURL, _ := url.Parse(fmt.Sprintf("http://localhost-%d", i))
		_, err := service.CreateArticle(ctx, testArticle(testURL.String()), testURL, nil, "")
		assert.NoError(t, err)
	}

	pageToken := ""
	for i := 0; i < 10; i++ {
		response, err := service.ListArticles(ctx, &ListArticlesRequest{
			PageSize:  5,
			PageToken: pageToken,
		})
		assert.NoError(t, err)
		assert.Equal(t, 5, len(response.Articles))

		if i != 9 {
			assert.NotEmpty(t, response.NextPageToken)
		} else {
			assert.Empty(t, response.NextPageToken)
		}

		pageToken = response.NextPageToken
	}
}
