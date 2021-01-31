package db

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type migration struct {
	Name  string
	Query string
}

func migrations(db *sql.DB) []*migration {
	binaryTypeName := "BLOB"
	if _, psql := db.Driver().(*pq.Driver); psql {
		binaryTypeName = "BYTEA"
	}

	return []*migration{
		{
			Name: "create users",
			Query: fmt.Sprintf(`
			CREATE TABLE users (
				id                 TEXT   NOT NULL,
				username           TEXT   NOT NULL,
				hash               %s     NOT NULL,
				created_epoch_utc  BIGINT NOT NULL,
				PRIMARY KEY (id),
				UNIQUE(username)
			)
			`, binaryTypeName),
		},
		{
			Name: "create jwt_keys",
			Query: fmt.Sprintf(`
			CREATE TABLE jwt_keys (
				id         TEXT NOT NULL,
				public_der %s   NOT NULL,
				PRIMARY KEY (id)
			)
			`, binaryTypeName),
		},
		{
			Name: "create operations",
			Query: `
			CREATE TABLE operations (
				id       TEXT    NOT NULL,
				user_id  TEXT    NOT NULL REFERENCES users(id),
				done     BOOLEAN NOT NULL,
				error    TEXT        NULL,
				response TEXT        NULL,
				PRIMARY KEY (id)
			)
			`,
		},
		{
			Name: "create feeds",
			Query: `
			CREATE TABLE feeds (
				id                TEXT   NOT NULL,
				url               TEXT   NOT NULL,
				title             TEXT   NOT NULL,
				created_epoch_utc BIGINT NOT NULL,
				updated_epoch_utc BIGINT     NULL,
				icon_url          TEXT       NULL,
				PRIMARY KEY (id),
				UNIQUE (url)
			)
			`,
		},
		{
			Name: "create users_feeds",
			Query: `
			CREATE TABLE users_feeds (
				user_id TEXT NOT NULL REFERENCES users(id),
				feed_id TEXT NOT NULL REFERENCES feeds(id),
				UNIQUE(user_id, feed_id)
			)
			`,
		},
		{
			Name: "create tags",
			Query: `
			CREATE TABLE tags (
				id                TEXT   NOT NULL,
				title             TEXT   NOT NULL,
				user_id           TEXT   NOT NULL REFERENCES users(id),
				created_epoch_utc BIGINT NOT NULL,
				PRIMARY KEY (id),
				UNIQUE (user_id, title)
			)
			`,
		},
		{
			Name: "create tags_feeds",
			Query: `
			CREATE TABLE tags_feeds (
				tag_id  TEXT NOT NULL REFERENCES tags(id),
				feed_id TEXT NOT NULL REFERENCES feeds(id),
				UNIQUE(tag_id, feed_id)
			)
			`,
		},
	}
}
