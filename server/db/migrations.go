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
	}
}
