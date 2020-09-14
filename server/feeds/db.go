package feeds

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
)

type dbFeed struct {
	Feed        *Feed
	LastFetched int64
}

type feedsDB struct {
	db *sql.DB
}

func newDB(sqldb *sql.DB) *feedsDB {
	return &feedsDB{
		db: sqldb,
	}
}

// Create adds a new articles to the database.
func (db *feedsDB) Create(ctx context.Context, feed *Feed) error {
	lastFetched, err := ptypes.Timestamp(feed.LastFetched)
	if err != nil {
		return fmt.Errorf("failed to convret last_fetched: %w", err)
	}

	_, createErr := db.db.ExecContext(ctx, `
	INSERT INTO feeds (
		id,
		user_id,
		url,
		title,
		last_fetched
	) VALUES (
		$1, $2, $3, $4, $5
	)
	`,
		feed.Id,
		feed.UserId,
		feed.Url,
		feed.Title,
		lastFetched.UnixNano(),
	)
	return createErr
}

// Get returns article by id.
func (db *feedsDB) Get(ctx context.Context, id string) (*Feed, error) {
	row := db.db.QueryRowContext(ctx, `
	SELECT
		id,
		user_id,
		url,
		title,
		last_fetched
	FROM
		feeds
	WHERE
		id = $1
	`, id)

	return db.scanRow(row)
}

type scannable interface {
	Scan(...interface{}) error
}

func (db *feedsDB) scanRow(row scannable) (*Feed, error) {
	feed := &dbFeed{
		Feed: &Feed{},
	}
	err := row.Scan(
		&feed.Feed.Id,
		&feed.Feed.UserId,
		&feed.Feed.Url,
		&feed.Feed.Title,
		&feed.LastFetched,
	)

	if err != nil {
		return nil, err
	}

	var convertTimeErr error
	feed.Feed.LastFetched, convertTimeErr = ptypes.TimestampProto(time.Unix(0, feed.LastFetched))
	if convertTimeErr != nil {
		return nil, fmt.Errorf("failed to convert create time: %w", convertTimeErr)
	}

	return feed.Feed, nil
}

func (db *feedsDB) Update(ctx context.Context, feed *Feed) error {
	lastFetched, err := ptypes.Timestamp(feed.LastFetched)
	if err != nil {
		return fmt.Errorf("failed to convret last_fetched: %w", err)
	}

	_, updateError := db.db.ExecContext(ctx, `
	UPDATE feeds
	SET
		last_fetched = $1
	WHERE
		id = $2
	`, lastFetched.UnixNano(), feed.Id)
	return updateError
}

// ListAll returns all articles.
func (db *feedsDB) ListAll(ctx context.Context) ([]*Feed, error) {
	rows, err := db.db.QueryContext(ctx, `
		SELECT
			id,
			user_id,
			url,
			title,
			last_fetched
		FROM
			feeds
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	feeds := []*Feed{}
	for rows.Next() {
		feed, err := db.scanRow(rows)
		if err != nil {
			return nil, err
		}

		feeds = append(feeds, feed)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return feeds, nil
}

// List returns articles for a user.
func (db *feedsDB) List(ctx context.Context, request *ListFeedsRequest) ([]*Feed, error) {
	from, err := request.FromID()
	if err != nil {
		return nil, err
	}

	rows, err := db.db.QueryContext(ctx, `
		SELECT
			id,
			user_id,
			url,
			title,
			last_fetched
		FROM
			feeds
		WHERE
			user_id = $1
			AND id >= $2
		ORDER BY id ASC
		LIMIT $3
	`, request.UserId, from, request.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	feeds := []*Feed{}
	for rows.Next() {
		feed, err := db.scanRow(rows)
		if err != nil {
			return nil, err
		}

		feeds = append(feeds, feed)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return feeds, nil
}
