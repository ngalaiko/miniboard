package items

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type logger interface {
	Error(string, ...interface{})
}

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

// Create creates a item in the database.
func (d *database) Create(ctx context.Context, item *Item) error {
	_, err := d.db.ExecContext(ctx, `
		INSERT INTO items (
			id,
			user_id,
			url,
			title,
			subscription_id,
			created_epoch_utc
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)`, item.ID, item.UserID, item.URL, item.Title,
		item.SubscriptionID, item.Created.UTC().UnixNano(),
	)
	return err
}

// Get returns a item from the db with the given id and user id.
func (d *database) Get(ctx context.Context, userID string, id string) (*Item, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT
		id,
		user_id,
		url,
		title,
		subscription_id,
		created_epoch_utc
	FROM
		items
	WHERE
		user_id = $1
		AND id = $2
	`, userID, id)

	return d.scanRow(row)
}

// List returns a list of items from the database.
func (d *database) List(ctx context.Context,
	userID string,
	limit int,
	createdLT *time.Time,
	subscriptionID *string,
) ([]*Item, error) {
	query := &strings.Builder{}
	query.WriteString(`
	SELECT
		id,
		user_id,
		url,
		title,
		subscription_id,
		created_epoch_utc
	FROM
		items
	WHERE
		user_id = $1
	`)
	args := []interface{}{userID}

	if subscriptionID != nil {
		args = append(args, *subscriptionID)
		query.WriteString(fmt.Sprintf(`
		AND subscription_id = $%d
		`, len(args)))
	}

	if createdLT != nil {
		args = append(args, createdLT.UnixNano())
		query.WriteString(fmt.Sprintf(`
		AND created_epoch_utc < $%d
		`, len(args)))
	}

	args = append(args, limit)
	query.WriteString(fmt.Sprintf(`
	ORDER BY
		items.created_epoch_utc DESC
	LIMIT $%d`, len(args)))

	rows, err := d.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}

	items := []*Item{}
	for rows.Next() {
		item, err := d.scanRow(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

type scannable interface {
	Scan(...interface{}) error
}

func (d *database) scanRow(row scannable) (*Item, error) {
	item := &Item{}
	var createdEpoch int64
	if err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.URL,
		&item.Title,
		&item.SubscriptionID,
		&createdEpoch,
	); err != nil {
		return nil, err
	}

	item.Created = time.Unix(0, createdEpoch).Round(time.Nanosecond)

	return item, nil
}
