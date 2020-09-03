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
				feed_id        TEXT    NULL,
				PRIMARY KEY (id),
				UNIQUE (user_id, content_sha256)
			)
			`,
		},
	}
}
