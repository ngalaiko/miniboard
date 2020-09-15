package db

type migration struct {
	Name  string
	Query string
}

func migrations() []*migration {
	return []*migration{
		{
			Name: "create public_keys",
			Query: `
			CREATE TABLE public_keys (
				id         TEXT NOT NULL,
				der_base64 TEXT NOT NULL,
				PRIMARY KEY (id)
			)
			`,
		},
		{
			Name: "create feeds table",
			Query: `
			CREATE TABLE feeds (
				id           TEXT   NOT NULL,
				user_id      TEXT   NOT NULL,
				url          TEXT   NOT NULL,
				title        TEXT   NOT NULL,
				last_fetched BIGINT NOT NULL,
				PRIMARY KEY (id),
				UNIQUE (user_id, url)
			)
			`,
		},
		{
			Name: "create articles table",
			Query: `
			CREATE TABLE articles (
				id             TEXT    NOT NULL,
				user_id        TEXT    NOT NULL,
				url            TEXT    NOT NULL,
				title          TEXT    NOT NULL,
				create_time    BIGINT  NOT NULL,
				content        TEXT    NOT NULL,
				content_sha256 TEXT    NOT NULL,
				is_read        BOOLEAN NOT NULL,
				feed_id        TEXT    NULL     REFERENCES feeds(id),
				PRIMARY KEY (id),
				UNIQUE (user_id, content_sha256)
			)
			`,
		},
		{
			Name: "create operations table",
			Query: `
			CREATE TABLE operations (
				id       TEXT    NOT NULL,
				user_id  TEXT    NOT NULL,
				done     BOOLEAN NOT NULL DEFAULT false,
				error    TEXT        NULL,
				response TEXT        NULL,
				metadata TEXT    NOT NULL,
				PRIMARY KEY (id)
			)
			`,
		},
	}
}
