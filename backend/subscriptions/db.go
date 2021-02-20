package subscriptions

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
		subscriptions.id,
		users_subscriptions.user_id,
		subscriptions.url,
		subscriptions.title,
		subscriptions.created_epoch_utc,
		subscriptions.updated_epoch_utc,
		subscriptions.icon_url,
		%s(tags.id, ',')
	`, groupFunc)
}

// Create creates a subscription in the database.
func (d *database) Create(ctx context.Context, subscription *UserSubscription) error {
	var updatedEpoch *int64
	if subscription.Updated != nil {
		updatedEpoch = new(int64)
		*updatedEpoch = subscription.Updated.UTC().UnixNano()
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	onError := func(tx *sql.Tx, err error) error {
		if rollbackError := tx.Rollback(); rollbackError != nil {
			d.logger.Error("failed to rollback transaction when creating subscription: %s", err)
		}
		return err
	}

	existingSubscription, err := getByURL(ctx, tx, subscription.URL)
	switch err {
	case sql.ErrNoRows:
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO subscriptions (
			id,
			url,
			title,
			created_epoch_utc,
			updated_epoch_utc,
			icon_url
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)`, subscription.ID, subscription.URL, subscription.Title,
			subscription.Created.UTC().UnixNano(),
			updatedEpoch, subscription.IconURL,
		); err != nil {
			return onError(tx, err)
		}
		existingSubscription = &subscription.Subscription
		fallthrough
	case nil:
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO users_subscriptions (
			user_id, subscription_id
		) VALUES (
			$1, $2
		)`, subscription.UserID, existingSubscription.ID,
		); err != nil {
			return onError(tx, err)
		}

		for _, tagID := range subscription.TagIDs {
			if _, err := tx.ExecContext(ctx, `
			INSERT INTO tags_subscriptions (
				tag_id, subscription_id
			) VALUES (
				$1, $2
			)`, tagID, existingSubscription.ID,
			); err != nil {
				return onError(tx, err)
			}
		}

		return tx.Commit()
	default:
		return onError(tx, err)
	}
}

func getByURL(ctx context.Context, tx *sql.Tx, url string) (*Subscription, error) {
	row := tx.QueryRowContext(ctx, `
	SELECT
		subscriptions.id,
		subscriptions.url,
		subscriptions.title,
		subscriptions.created_epoch_utc,
		subscriptions.updated_epoch_utc,
		subscriptions.icon_url
	FROM
		subscriptions
	WHERE
		subscriptions.url = $1
	`, url)

	return scanSubscriptionRow(row)
}

// GetByID returns a subscription from the db with the given id.
func (d *database) GetByID(ctx context.Context, id string) (*Subscription, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT
		subscriptions.id,
		subscriptions.url,
		subscriptions.title,
		subscriptions.created_epoch_utc,
		subscriptions.updated_epoch_utc,
		subscriptions.icon_url
	FROM
		subscriptions
	WHERE
		id = $1
	`, id)

	return scanSubscriptionRow(row)
}

// Get returns a subscription from the db with the given url and user id.
func (d *database) GetByURL(ctx context.Context, userID string, url string) (*UserSubscription, error) {
	row := d.db.QueryRowContext(ctx, fmt.Sprintf(`
	SELECT
		%s
	FROM
		subscriptions
			JOIN users_subscriptions ON subscriptions.id = users_subscriptions.subscription_id AND users_subscriptions.user_id = $1
			LEFT JOIN tags_subscriptions ON subscriptions.id = tags_subscriptions.subscription_id
			LEFT JOIN tags ON tags.id = tags_subscriptions.tag_id AND tags.user_id = users_subscriptions.user_id
	WHERE
		subscriptions.url = $2
	GROUP BY
		subscriptions.id,
		users_subscriptions.user_id,
		subscriptions.url,
		subscriptions.title,
		subscriptions.created_epoch_utc,
		subscriptions.updated_epoch_utc,
		subscriptions.icon_url
	`, sqlFields(d.db)), userID, url)

	return scanUserSubscriptionRow(row)
}

// Get returns a subscription from the db with the given id and user id.
func (d *database) Get(ctx context.Context, userID string, id string) (*UserSubscription, error) {
	row := d.db.QueryRowContext(ctx, fmt.Sprintf(`
	SELECT
		%s
	FROM
		subscriptions
			JOIN users_subscriptions ON subscriptions.id = users_subscriptions.subscription_id AND users_subscriptions.user_id = $1
			LEFT JOIN tags_subscriptions ON subscriptions.id = tags_subscriptions.subscription_id
			LEFT JOIN tags ON tags.id = tags_subscriptions.tag_id AND tags.user_id = users_subscriptions.user_id
	WHERE
		subscriptions.id = $2
	GROUP BY
		subscriptions.id,
		users_subscriptions.user_id,
		subscriptions.url,
		subscriptions.title,
		subscriptions.created_epoch_utc,
		subscriptions.updated_epoch_utc,
		subscriptions.icon_url
	`, sqlFields(d.db)), userID, id)

	return scanUserSubscriptionRow(row)
}

// List returns a list of subscriptions from the database.
func (d *database) List(ctx context.Context,
	userID string,
	limit int,
	createdLT *time.Time,
) ([]*UserSubscription, error) {
	query := &strings.Builder{}
	query.WriteString(fmt.Sprintf(`
	SELECT
		%s
	FROM
		subscriptions
			JOIN users_subscriptions ON subscriptions.id = users_subscriptions.subscription_id AND users_subscriptions.user_id = $1
			LEFT JOIN tags_subscriptions ON subscriptions.id = tags_subscriptions.subscription_id
			LEFT JOIN tags ON tags.id = tags_subscriptions.tag_id AND tags.user_id = users_subscriptions.user_id
	`, sqlFields(d.db)))

	args := []interface{}{userID}

	if createdLT != nil {
		args = append(args, createdLT.UnixNano())
		query.WriteString(fmt.Sprintf(`
	WHERE
		subscriptions.created_epoch_utc < $%d
		`, len(args)))
	}
	args = append(args, limit)

	query.WriteString(fmt.Sprintf(`
	GROUP BY
		subscriptions.id,
		users_subscriptions.user_id,
		subscriptions.url,
		subscriptions.title,
		subscriptions.created_epoch_utc,
		subscriptions.updated_epoch_utc,
		subscriptions.icon_url
	ORDER BY
		subscriptions.created_epoch_utc DESC
	LIMIT $%d`, len(args)))

	rows, err := d.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}

	subscriptions := []*UserSubscription{}
	for rows.Next() {
		subscription, err := scanUserSubscriptionRow(rows)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// ListAll returns a list of all subscriptions from the database.
func (d *database) ListAll(ctx context.Context) ([]*Subscription, error) {
	query := `
	SELECT
		subscriptions.id,
		subscriptions.url,
		subscriptions.title,
		subscriptions.created_epoch_utc,
		subscriptions.updated_epoch_utc,
		subscriptions.icon_url
	FROM
		subscriptions`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	ss := []*Subscription{}
	for rows.Next() {
		s, err := scanSubscriptionRow(rows)
		if err != nil {
		}
		ss = append(ss, s)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ss, nil
}

type scannable interface {
	Scan(...interface{}) error
}

func scanSubscriptionRow(row scannable) (*Subscription, error) {
	subscription := &Subscription{}
	var createdEpoch int64
	var updatedEpoch *int64
	if err := row.Scan(
		&subscription.ID,
		&subscription.URL,
		&subscription.Title,
		&createdEpoch,
		&updatedEpoch,
		&subscription.IconURL,
	); err != nil {
		return nil, err
	}

	subscription.Created = time.Unix(0, createdEpoch).Round(time.Nanosecond)

	if updatedEpoch != nil {
		subscription.Updated = new(time.Time)
		*subscription.Updated = time.Unix(0, *updatedEpoch).Round(time.Nanosecond)
	}

	return subscription, nil
}

func scanUserSubscriptionRow(row scannable) (*UserSubscription, error) {
	subscription := &UserSubscription{}
	var createdEpoch int64
	var updatedEpoch *int64
	if err := row.Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.URL,
		&subscription.Title,
		&createdEpoch,
		&updatedEpoch,
		&subscription.IconURL,
		&subscription.TagIDs,
	); err != nil {
		return nil, err
	}

	subscription.Created = time.Unix(0, createdEpoch).Round(time.Nanosecond)

	if updatedEpoch != nil {
		subscription.Updated = new(time.Time)
		*subscription.Updated = time.Unix(0, *updatedEpoch).Round(time.Nanosecond)
	}

	return subscription, nil
}
