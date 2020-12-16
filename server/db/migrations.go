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
				id   TEXT NOT NULL,
				hash BLOB NOT NULL,
				PRIMARY KEY (id)
			)
			`,
		},
		{
			Name: "create public_keys",
			Query: `
			CREATE TABLE public_keys (
				id  TEXT NOT NULL,
				der BLOB NOT NULL,
				PRIMARY KEY (id)
			)
			`,
		},
	}
}
