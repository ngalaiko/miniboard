package feeds

import (
	"context"
	"database/sql"
	"net/url"
	"time"
)

type database struct {
	db *sql.DB
}

func newDB(sqldb *sql.DB) *database {
	return &database{
		db: sqldb,
	}
}

// Create creates a feed in the database.
func (d *database) Create(ctx context.Context, feed *Feed) error {
	var updatedEpoch *int64
	if feed.Updated != nil {
		updatedEpoch = new(int64)
		*updatedEpoch = feed.Updated.UTC().UnixNano()
	}

	_, err := d.db.ExecContext(ctx, `
	INSERT INTO feeds (
		id,
		user_id,
		url,
		title,
		created_epoch,
		updated_epoch
	) VALUES (
		$1, $2, $3, $4, $5, $6
	)`, feed.ID, feed.UserID, feed.URL.String(), feed.Title,
		feed.Created.UTC().UnixNano(),
		updatedEpoch,
	)

	return err
}

// Get returns a feed from the db with the given id and user id.
func (d *database) Get(ctx context.Context, userID string, id string) (*Feed, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT
		id,
		user_id,
		url,
		title,
		created_epoch,
		updated_epoch
	FROM
		feeds
	WHERE
		id = $1
		AND user_id = $2
	`, id, userID)

	return d.scanRow(row)
}

type scannable interface {
	Scan(...interface{}) error
}

func (d *database) scanRow(row scannable) (*Feed, error) {
	feed := &Feed{}
	var rawURL string
	var createdEpoch int64
	var updatedEpoch *int64
	if err := row.Scan(
		&feed.ID,
		&feed.UserID,
		&rawURL,
		&feed.Title,
		&createdEpoch,
		&updatedEpoch,
	); err != nil {
		return nil, err
	}

	var err error
	feed.URL, err = url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}

	feed.Created = time.Unix(0, createdEpoch).Round(time.Nanosecond)

	if updatedEpoch != nil {
		feed.Updated = new(time.Time)
		*feed.Updated = time.Unix(0, *updatedEpoch).Round(time.Nanosecond)
	}

	return feed, nil
}
