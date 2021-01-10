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
				id       TEXT NOT NULL,
				username TEXT NOT NULL,
				hash     BLOB NOT NULL,
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
				id       TEXT NOT NULL,
				user_id  TEXT NOT NULL,
				done     BOOLEAN NOT NULL,
				error    TEXT NULL,
				response TEXT NULL,
				PRIMARY KEY (id)
			)
			`,
		},
		{
			Name: "create feeds",
			Query: `
			CREATE TABLE feeds (
				id            TEXT   NOT NULL,
				user_id       TEXT   NOT NULL,
				url           TEXT   NOT NULL,
				title         TEXT   NOT NULL,
				created_epoch BIGINT NOT NULL,
				updated_epoch BIGINT     NULL,
				PRIMARY KEY (id)
			)
			`,
		},
	}
}
