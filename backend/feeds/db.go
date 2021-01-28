package feeds

import (
	"context"
	"database/sql"
	"strings"
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

	existingFeed, err := d.getByURL(ctx, feed.URL)
	switch err {
	case sql.ErrNoRows:
		if _, err := d.db.ExecContext(ctx, `
		INSERT INTO feeds (
			id,
			url,
			title,
			created_epoch,
			updated_epoch,
			icon_url
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)`, feed.ID, feed.URL, feed.Title,
			feed.Created.UTC().UnixNano(),
			updatedEpoch, feed.IconURL,
		); err != nil {
			return err
		}
		existingFeed = feed
		fallthrough
	case nil:
		if _, err := d.db.ExecContext(ctx, `
		INSERT INTO users_feeds (
			user_id, feed_id
		) VALUES (
			$1, $2
		)`, feed.UserID, existingFeed.ID,
		); err != nil {
			return err
		}
		return nil
	default:
		return err
	}
}

func (d *database) getByURL(ctx context.Context, url string) (*Feed, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT
		feeds.id,
		'',
		feeds.url,
		feeds.title,
		feeds.created_epoch,
		feeds.updated_epoch,
		feeds.icon_url
	FROM
		feeds
	WHERE
		feeds.url = $1
	`, url)

	return d.scanRow(row)
}

// Get returns a feed from the db with the given url and user id.
func (d *database) GetByURL(ctx context.Context, userID string, url string) (*Feed, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT
		feeds.id,
		users_feeds.user_id,
		feeds.url,
		feeds.title,
		feeds.created_epoch,
		feeds.updated_epoch,
		feeds.icon_url
	FROM
		feeds JOIN users_feeds ON feeds.id = users_feeds.feed_id AND users_feeds.user_id = $1
	WHERE
		feeds.url = $2
	`, userID, url)

	return d.scanRow(row)
}

// Get returns a feed from the db with the given id and user id.
func (d *database) Get(ctx context.Context, userID string, id string) (*Feed, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT
		feeds.id,
		users_feeds.user_id,
		feeds.url,
		feeds.title,
		feeds.created_epoch,
		feeds.updated_epoch,
		feeds.icon_url
	FROM
		feeds JOIN users_feeds ON feeds.id = users_feeds.feed_id AND users_feeds.user_id = $1
	WHERE
		feeds.id = $2
	`, userID, id)

	return d.scanRow(row)
}

// List returns a list of feeds from the database.
func (d *database) List(ctx context.Context, userID string, limit int, createdLT *time.Time) ([]*Feed, error) {
	query := &strings.Builder{}
	query.WriteString(`
	SELECT
		feeds.id,
		users_feeds.user_id,
		feeds.url,
		feeds.title,
		feeds.created_epoch,
		feeds.updated_epoch,
		feeds.icon_url
	FROM
		feeds JOIN users_feeds ON feeds.id = users_feeds.feed_id AND users_feeds.user_id = $1
	`)

	args := []interface{}{userID}
	if createdLT != nil {
		query.WriteString(`WHERE feeds.created_epoch < $2 ORDER BY feeds.created_epoch DESC LIMIT $3`)
		args = append(args, createdLT.UnixNano(), limit)
	} else {
		query.WriteString(`ORDER BY feeds.created_epoch DESC LIMIT $2`)
		args = append(args, limit)
	}

	rows, err := d.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}

	feeds := []*Feed{}
	for rows.Next() {
		feed, err := d.scanRow(rows)
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

type scannable interface {
	Scan(...interface{}) error
}

func (d *database) scanRow(row scannable) (*Feed, error) {
	feed := &Feed{}
	var createdEpoch int64
	var updatedEpoch *int64
	if err := row.Scan(
		&feed.ID,
		&feed.UserID,
		&feed.URL,
		&feed.Title,
		&createdEpoch,
		&updatedEpoch,
		&feed.IconURL,
	); err != nil {
		return nil, err
	}

	feed.Created = time.Unix(0, createdEpoch).Round(time.Nanosecond)

	if updatedEpoch != nil {
		feed.Updated = new(time.Time)
		*feed.Updated = time.Unix(0, *updatedEpoch).Round(time.Nanosecond)
	}

	return feed, nil
}
