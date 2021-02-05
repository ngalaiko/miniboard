package feeds

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

type database struct {
	db     *sql.DB
	logger logger
}

func newDB(sqldb *sql.DB, logger logger) *database {
	return &database{
		db:     sqldb,
		logger: logger,
	}
}

func sqlFields(db *sql.DB) string {
	groupFunc := "GROUP_CONCAT"
	if _, ok := db.Driver().(*pq.Driver); ok {
		groupFunc = "STRING_AGG"
	}

	return fmt.Sprintf(`
		feeds.id,
		users_feeds.user_id,
		feeds.url,
		feeds.title,
		feeds.created_epoch_utc,
		feeds.updated_epoch_utc,
		feeds.icon_url,
		%s(tags_feeds.tag_id, ',')
	`, groupFunc)
}

// Create creates a feed in the database.
func (d *database) Create(ctx context.Context, feed *Feed) error {
	var updatedEpoch *int64
	if feed.Updated != nil {
		updatedEpoch = new(int64)
		*updatedEpoch = feed.Updated.UTC().UnixNano()
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	onError := func(tx *sql.Tx, err error) error {
		if rollbackError := tx.Rollback(); rollbackError != nil {
			d.logger.Error("failed to rollback transaction when creating feed: %s", err)
		}
		return err
	}

	existingFeed, err := d.getByURL(ctx, feed.URL)
	switch err {
	case sql.ErrNoRows:
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO feeds (
			id,
			url,
			title,
			created_epoch_utc,
			updated_epoch_utc,
			icon_url
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)`, feed.ID, feed.URL, feed.Title,
			feed.Created.UTC().UnixNano(),
			updatedEpoch, feed.IconURL,
		); err != nil {
			return onError(tx, err)
		}
		existingFeed = feed
		fallthrough
	case nil:
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO users_feeds (
			user_id, feed_id
		) VALUES (
			$1, $2
		)`, feed.UserID, existingFeed.ID,
		); err != nil {
			return onError(tx, err)
		}

		for _, tagID := range feed.TagIDs {
			if _, err := tx.ExecContext(ctx, `
			INSERT INTO tags_feeds (
				tag_id, feed_id
			) VALUES (
				$1, $2
			)`, tagID, existingFeed.ID,
			); err != nil {
				return onError(tx, err)
			}
		}

		return tx.Commit()
	default:
		return onError(tx, err)
	}
}

func (d *database) getByURL(ctx context.Context, url string) (*Feed, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT
		feeds.id,
		'',
		feeds.url,
		feeds.title,
		feeds.created_epoch_utc,
		feeds.updated_epoch_utc,
		feeds.icon_url,
		NULL
	FROM
		feeds
	WHERE
		feeds.url = $1
	`, url)

	return d.scanRow(row)
}

// Get returns a feed from the db with the given url and user id.
func (d *database) GetByURL(ctx context.Context, userID string, url string) (*Feed, error) {
	row := d.db.QueryRowContext(ctx, fmt.Sprintf(`
	SELECT
		%s
	FROM
		feeds
			JOIN users_feeds ON feeds.id = users_feeds.feed_id AND users_feeds.user_id = $1
			LEFT JOIN tags_feeds ON feeds.id = tags_feeds.feed_id
	WHERE
		feeds.url = $2
	GROUP BY
		feeds.id,
		users_feeds.user_id,
		feeds.url,
		feeds.title,
		feeds.created_epoch_utc,
		feeds.updated_epoch_utc,
		feeds.icon_url
	`, sqlFields(d.db)), userID, url)

	return d.scanRow(row)
}

// Get returns a feed from the db with the given id and user id.
func (d *database) Get(ctx context.Context, userID string, id string) (*Feed, error) {
	row := d.db.QueryRowContext(ctx, fmt.Sprintf(`
	SELECT
		%s
	FROM
		feeds
			JOIN users_feeds ON feeds.id = users_feeds.feed_id AND users_feeds.user_id = $1
			LEFT JOIN tags_feeds ON feeds.id = tags_feeds.feed_id
	WHERE
		feeds.id = $2
	GROUP BY
		feeds.id,
		users_feeds.user_id,
		feeds.url,
		feeds.title,
		feeds.created_epoch_utc,
		feeds.updated_epoch_utc,
		feeds.icon_url
	`, sqlFields(d.db)), userID, id)

	return d.scanRow(row)
}

// List returns a list of feeds from the database.
func (d *database) List(ctx context.Context,
	userID string,
	limit int,
	createdLT *time.Time,
) ([]*Feed, error) {
	query := &strings.Builder{}
	query.WriteString(fmt.Sprintf(`
	SELECT
		%s
	FROM
		feeds
			JOIN users_feeds ON feeds.id = users_feeds.feed_id AND users_feeds.user_id = $1
			LEFT JOIN tags_feeds ON feeds.id = tags_feeds.feed_id
	`, sqlFields(d.db)))

	args := []interface{}{userID}

	if createdLT != nil {
		args = append(args, createdLT.UnixNano())
		query.WriteString(fmt.Sprintf(`
	WHERE
		feeds.created_epoch_utc < $%d
		`, len(args)))
	}
	args = append(args, limit)

	query.WriteString(fmt.Sprintf(`
	GROUP BY
		feeds.id,
		users_feeds.user_id,
		feeds.url,
		feeds.title,
		feeds.created_epoch_utc,
		feeds.updated_epoch_utc,
		feeds.icon_url
	ORDER BY
		feeds.created_epoch_utc DESC
	LIMIT $%d`, len(args)))

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
		&feed.TagIDs,
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
