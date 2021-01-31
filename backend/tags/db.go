package tags

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

// Create creates a tag in the database.
func (d *database) Create(ctx context.Context, tag *Tag) error {
	if _, err := d.db.ExecContext(ctx, `
		INSERT INTO tags (
			id,
			user_id,
			title,
			created_epoch
		) VALUES (
			$1, $2, $3, $4
		)`, tag.ID, tag.UserID, tag.Title, tag.Created.UTC().UnixNano(),
	); err != nil {
		return err
	}
	return nil
}

// Get returns a tag from the db with the given id and user id.
func (d *database) GetByID(ctx context.Context, userID string, id string) (*Tag, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT
		id,
		user_id,
		title,
		created_epoch
	FROM
		tags
	WHERE
		id = $1
		AND user_id = $2
	`, id, userID)

	return d.scanRow(row)
}

// Get returns a tag from the db with the given title and user id.
func (d *database) GetByTitle(ctx context.Context, userID string, title string) (*Tag, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT
		id,
		user_id,
		title,
		created_epoch
	FROM
		tags
	WHERE
		user_id = $2
		AND title = $1
	`, userID, title)

	return d.scanRow(row)
}

// List returns a list of tags from the database.
func (d *database) List(ctx context.Context, userID string, limit int, createdLT *time.Time) ([]*Tag, error) {
	query := &strings.Builder{}
	query.WriteString(`
	SELECT
		id,
		user_id,
		title,
		created_epoch
	FROM
		tags
	WHERE
		user_id = $1
	`)

	args := []interface{}{userID}
	if createdLT != nil {
		query.WriteString(`AND created_epoch < $2 ORDER BY created_epoch DESC LIMIT $3`)
		args = append(args, createdLT.UnixNano(), limit)
	} else {
		query.WriteString(`ORDER BY created_epoch DESC LIMIT $2`)
		args = append(args, limit)
	}

	rows, err := d.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}

	tags := []*Tag{}
	for rows.Next() {
		tag, err := d.scanRow(rows)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

type scannable interface {
	Scan(...interface{}) error
}

func (d *database) scanRow(row scannable) (*Tag, error) {
	tag := &Tag{}
	var createdEpoch int64
	if err := row.Scan(
		&tag.ID,
		&tag.UserID,
		&tag.Title,
		&createdEpoch,
	); err != nil {
		return nil, err
	}

	tag.Created = time.Unix(0, createdEpoch).Round(time.Nanosecond)

	return tag, nil
}
