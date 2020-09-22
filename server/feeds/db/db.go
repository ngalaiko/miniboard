package db

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/ngalaiko/miniboard/server/genproto/feeds/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type dbFeed struct {
	Feed        *feeds.Feed
	LastFetched int64
}

// DB allows to access feeds database.
type DB struct {
	db *sql.DB
}

// New returns new feeds database.
func New(sqldb *sql.DB) *DB {
	return &DB{
		db: sqldb,
	}
}

// Create adds a new articles to the database.
func (db *DB) Create(ctx context.Context, feed *feeds.Feed) error {
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
func (db *DB) Get(ctx context.Context, id string, userID string) (*feeds.Feed, error) {
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
		AND user_id = $2
	`, id, userID)

	return db.scanRow(row)
}

type scannable interface {
	Scan(...interface{}) error
}

func (db *DB) scanRow(row scannable) (*feeds.Feed, error) {
	feed := &dbFeed{
		Feed: &feeds.Feed{},
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

// Upates updates a feed.
func (db *DB) Update(ctx context.Context, feed *feeds.Feed, userID string) error {
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
		AND user_id = $3
	`, lastFetched.UnixNano(), feed.Id, userID)
	return updateError
}

// ListAll returns all articles.
func (db *DB) ListAll(ctx context.Context) ([]*feeds.Feed, error) {
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
	feeds := []*feeds.Feed{}
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
func (db *DB) List(ctx context.Context, userID string, request *feeds.ListFeedsRequest) ([]*feeds.Feed, error) {
	from, err := getID(request)
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
	`, userID, from, request.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	feeds := []*feeds.Feed{}
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

func getID(request *feeds.ListFeedsRequest) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(request.PageToken)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "invalid page token")
	}
	return string(decoded), nil
}
