package db

type migration struct {
	Name  string
	Query string
}

func migrations() []*migration {
	return []*migration{
		{
			Name: "create users",
			Query: `
			CREATE TABLE users (
				id            TEXT   NOT NULL,
				username      TEXT   NOT NULL,
				hash          BLOC   NOT NULL,
				created_epoch BIGINT NOT NULL,
				PRIMARY KEY (id),
				UNIQUE(username)
			)
			`,
		},
		{
			Name: "create jwt_keys",
			Query: `
			CREATE TABLE jwt_keys (
				id         TEXT NOT NULL,
				public_der BLOB NOT NULL,
				PRIMARY KEY (id)
			)
			`,
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
				icon_url      TEXT       NULL,
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
		{
			Name: "create items_subscription_id_created_epoch_ix",
			Query: `
			CREATE INDEX items_subscription_id_created_epoch_ix
			ON items (subscription_id, created_epoch)
			`,
		},
		{
			Name: "create subscriptions_url_ix",
			Query: `
			CREATE INDEX subscriptions_url_ix
			ON subscriptions (url)
			`,
		},
		{
			Name: "create tags_user_id_title_ix",
			Query: `
			CREATE INDEX tags_user_id_title_ix
			ON tags (user_id, title)
			`,
		},
		{
			Name: "create tags_user_id_created_epoch_ix",
			Query: `
			CREATE INDEX tags_user_id_created_epoch_ix
			ON tags (user_id, created_epoch)
			`,
		},
	}
}
