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
				id            TEXT   NOT NULL,
				username      TEXT   NOT NULL,
				hash          %s     NOT NULL,
				created_epoch BIGINT NOT NULL,
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
			Name: "create subscriptions",
			Query: `
			CREATE TABLE subscriptions (
				id            TEXT   NOT NULL,
				url           TEXT   NOT NULL,
				title         TEXT   NOT NULL,
				created_epoch BIGINT NOT NULL,
				updated_epoch BIGINT     NULL,
				icon_url          TEXT       NULL,
				PRIMARY KEY (id),
				UNIQUE (url)
			)
			`,
		},
		{
			Name: "create users_subscriptions",
			Query: `
			CREATE TABLE users_subscriptions (
				user_id TEXT NOT NULL REFERENCES users(id),
				subscription_id TEXT NOT NULL REFERENCES subscriptions(id),
				UNIQUE(user_id, subscription_id)
			)
			`,
		},
		{
			Name: "create tags",
			Query: `
			CREATE TABLE tags (
				id            TEXT   NOT NULL,
				title         TEXT   NOT NULL,
				user_id       TEXT   NOT NULL REFERENCES users(id),
				created_epoch BIGINT NOT NULL,
				PRIMARY KEY (id),
				UNIQUE (user_id, title)
			)
			`,
		},
		{
			Name: "create tags_subscriptions",
			Query: `
			CREATE TABLE tags_subscriptions (
				tag_id  TEXT NOT NULL REFERENCES tags(id),
				subscription_id TEXT NOT NULL REFERENCES subscriptions(id),
				UNIQUE(tag_id, subscription_id)
			)
			`,
		},
		{
			Name: "create items",
			Query: `
			CREATE TABLE items (
				id              TEXT   NOT NULL,
				url             TEXT   NOT NULL,
				title           TEXT   NOT NULL,
				subscription_id TEXT   NOT NULL REFERENCES subscriptions(id),
				created_epoch   BIGINT NULL,
				summary         TEXT   NULL,
				PRIMARY KEY (id),
				UNIQUE (subscription_id, url)
			)
			`,
		},
	}
}
