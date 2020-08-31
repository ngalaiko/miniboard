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
	}
}
