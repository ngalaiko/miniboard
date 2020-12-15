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
				id         TEXT NOT NULL,
				PRIMARY KEY (id)
			)
			`,
		},
	}
}
