package articles

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/ngalaiko/miniboard/server/actor"
)

type dbArticle struct {
	Article         *Article
	CreateTimestamp int64
	FeedID          *string
}

type articlesDB struct {
	db *sql.DB
}

func newDB(sqldb *sql.DB) *articlesDB {
	return &articlesDB{
		db: sqldb,
	}
}

// Create adds a new articles to the database.
func (db *articlesDB) Create(ctx context.Context, article *Article) error {
	createTime, err := ptypes.Timestamp(article.CreateTime)
	if err != nil {
		return fmt.Errorf("failed to convret create_time: %w", err)
	}

	var feedID *string
	if article.FeedId != nil {
		feedID = new(string)
		*feedID = article.FeedId.GetValue()
	}

	_, createErr := db.db.ExecContext(ctx, `
	INSERT INTO articles (
		id,
		user_id,
		url,
		title,
		create_time,
		content,
		content_sha256,
		is_read,
		feed_id
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9
	)
	`,
		article.Id,
		article.UserId,
		article.Url,
		article.Title,
		createTime.UnixNano(),
		article.Content,
		article.ContentSha256,
		article.IsRead,
		feedID,
	)
	return createErr
}

// Get returns article by id.
func (db *articlesDB) Get(ctx context.Context, id string, userID string) (*Article, error) {
	row := db.db.QueryRowContext(ctx, `
	SELECT
		id,
		user_id,
		url,
		title,
		create_time,
		content,
		content_sha256,
		is_read,
		feed_id
	FROM
		articles
	WHERE
		id = $1
		AND user_id = $2
	`, id, userID)

	return db.scanRow(row)
}

// GetByUserIDUrl returns article by userId and sha256 sum.
func (db *articlesDB) GetByUserIDUrl(ctx context.Context, userID string, url string) (*Article, error) {
	row := db.db.QueryRowContext(ctx, `
	SELECT
		id,
		user_id,
		url,
		title,
		create_time,
		content,
		content_sha256,
		is_read,
		feed_id
	FROM
		articles
	WHERE
		user_id = $1
		AND url = $2
	`, userID, url)

	return db.scanRow(row)
}

type scannable interface {
	Scan(...interface{}) error
}

func (db *articlesDB) scanRow(row scannable) (*Article, error) {
	article := &dbArticle{
		Article: &Article{},
	}
	err := row.Scan(
		&article.Article.Id,
		&article.Article.UserId,
		&article.Article.Url,
		&article.Article.Title,
		&article.CreateTimestamp,
		&article.Article.Content,
		&article.Article.ContentSha256,
		&article.Article.IsRead,
		&article.FeedID,
	)

	if err != nil {
		return nil, err
	}

	var convertTimeErr error
	article.Article.CreateTime, convertTimeErr = ptypes.TimestampProto(time.Unix(0, article.CreateTimestamp))
	if convertTimeErr != nil {
		return nil, fmt.Errorf("failed to convert create time: %w", convertTimeErr)
	}

	if article.FeedID != nil {
		article.Article.FeedId = &wrappers.StringValue{
			Value: *article.FeedID,
		}
	}

	return article.Article, nil
}

// Delete deletes an article by id.
func (db *articlesDB) Delete(ctx context.Context, id string, userID string) error {
	_, err := db.db.ExecContext(ctx, `
	DELETE FROM articles
	WHERE
		id = $1
		AND user_id = $2
	`, id, userID)
	return err
}

func (db *articlesDB) Update(ctx context.Context, article *Article, userID string) error {
	_, err := db.db.ExecContext(ctx, `
	UPDATE articles
	SET
		is_read = $1,
		content = $2,
		content_sha256 = $3
	WHERE
		id = $4
		AND user_id = $5
	`, article.IsRead, article.Content, article.ContentSha256, article.Id, userID)
	return err
}

// List returns all articles.
func (db *articlesDB) List(ctx context.Context, request *ListArticlesRequest) ([]*Article, error) {
	a, _ := actor.FromContext(ctx)

	q := &strings.Builder{}
	q.WriteString(`
		SELECT
			id,
			user_id,
			url,
			title,
			create_time,
			content,
			content_sha256,
			is_read,
			feed_id
		FROM articles
		WHERE
			user_id = $1
			AND id >= $2
	`)

	from, err := request.FromID()
	if err != nil {
		return nil, err
	}

	args := []interface{}{
		a.ID, from,
	}

	if request.IsReadEq != nil {
		args = append(args, request.GetIsReadEq().GetValue())
		q.WriteString(fmt.Sprintf(`
			AND is_read = $%d
		`, len(args)))
	}

	if request.GetFeedIdEq() != nil {
		args = append(args, request.GetFeedIdEq().GetValue())
		q.WriteString(fmt.Sprintf(`
			AND feed_id = $%d
		`, len(args)))
	}

	if request.GetTitleContains() != nil {
		args = append(args, fmt.Sprintf("%%%s%%", request.GetTitleContains().GetValue()))
		q.WriteString(fmt.Sprintf(`
			AND title LIKE $%d
		`, len(args)))
	}

	args = append(args, request.PageSize)
	q.WriteString(fmt.Sprintf(`
		ORDER BY id ASC
		LIMIT $%d
	`, len(args)))

	rows, err := db.db.QueryContext(ctx, q.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	articles := []*Article{}
	for rows.Next() {
		article, err := db.scanRow(rows)
		if err != nil {
			return nil, err
		}
		article.Content = nil
		articles = append(articles, article)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return articles, nil
}
