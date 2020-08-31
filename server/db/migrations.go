package db

type migration struct {
	Name  string
	Query string
}

func migrations() []*migration {
	return []*migration{
		{
			Name: "create jwt_public_keys",
			Query: `
			CREATE TABLE jwt_public_keys (
				id             TEXT NOT NULL,
				public_key_der BLOB NOT NULL,
				PRIMARY KEY (id)
			)
			`,
		},
	}
}
