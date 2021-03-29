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
	var created *int64
	if item.Created != nil {
		nano := item.Created.UnixNano()
		created = &nano
	}

	_, err := d.db.ExecContext(ctx, `
		INSERT INTO items (
			id,
			url,
			title,
			subscription_id,
			created_epoch,
			summary
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)`, item.ID, item.URL, item.Title,
		item.SubscriptionID, created, item.Summary,
	)
	return err
}

// GetByURL returns a item from the db with the given url and user id.
func (d *database) GetByURL(ctx context.Context, url string) (*Item, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT
		items.id,
		items.url,
		items.title,
		items.subscription_id,
		items.created_epoch,
		items.summary
	FROM
		items
	WHERE
		items.url = $1
	`, url)

	return d.scanItemRow(row)
}

// Get returns a item from the db with the given id and user id.
func (d *database) Get(ctx context.Context, userID string, id string) (*UserItem, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT
		items.id,
		users_subscriptions.user_id,
		items.url,
		items.title,
		items.subscription_id,
		items.created_epoch,
		items.summary
	FROM
		items
			JOIN users_subscriptions on users_subscriptions.subscription_id = items.subscription_id AND users_subscriptions.user_id = $1
	WHERE
		id = $2
	`, userID, id)

	return d.scanUserItemRow(row)
}

// List returns a list of items from the database.
func (d *database) List(ctx context.Context,
	userID string,
	limit int,
	createdLT *time.Time,
	subscriptionID *string,
	tagID *string,
) ([]*UserItem, error) {
	query := &strings.Builder{}
	query.WriteString(`
	SELECT
		items.id,
		users_subscriptions.user_id,
		items.url,
		items.title,
		items.subscription_id,
		items.created_epoch,
		items.summary
	FROM
		items
			JOIN users_subscriptions on users_subscriptions.subscription_id = items.subscription_id AND users_subscriptions.user_id = $1
	`)
	args := []interface{}{userID}

	if tagID != nil {
		args = append(args, *tagID)
		query.WriteString(fmt.Sprintf(`
			JOIN tags_subscriptions on tags_subscriptions.subscription_id = users_subscriptions.subscription_id AND tags_subscriptions.tag_id = $%d
		`, len(args)))
	}

	if subscriptionID != nil {
		args = append(args, *subscriptionID)
		query.WriteString(fmt.Sprintf(`
		AND items.subscription_id = $%d
		`, len(args)))
	}

	if createdLT != nil {
		args = append(args, createdLT.UnixNano())
		query.WriteString(fmt.Sprintf(`
		AND items.created_epoch < $%d
		`, len(args)))
	}

	args = append(args, limit)
	query.WriteString(fmt.Sprintf(`
	ORDER BY
		items.created_epoch DESC
	LIMIT $%d`, len(args)))

	rows, err := d.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}

	items := []*UserItem{}
	for rows.Next() {
		item, err := d.scanUserItemRow(rows)
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

func (d *database) scanItemRow(row scannable) (*Item, error) {
	item := &Item{}
	var createdEpoch *int64
	if err := row.Scan(
		&item.ID,
		&item.URL,
		&item.Title,
		&item.SubscriptionID,
		&createdEpoch,
		&item.Summary,
	); err != nil {
		return nil, err
	}

	if createdEpoch != nil {
		item.Created = new(time.Time)
		*item.Created = time.Unix(0, *createdEpoch).Round(time.Nanosecond)
	}

	return item, nil
}

func (d *database) scanUserItemRow(row scannable) (*UserItem, error) {
	item := &UserItem{}
	var createdEpoch *int64
	if err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.URL,
		&item.Title,
		&item.SubscriptionID,
		&createdEpoch,
		&item.Summary,
	); err != nil {
		return nil, err
	}

	if createdEpoch != nil {
		item.Created = new(time.Time)
		*item.Created = time.Unix(0, *createdEpoch).Round(time.Nanosecond)
	}
	return item, nil
}
