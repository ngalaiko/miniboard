package jwt

// Token contains information about an authorized user.
type Token struct {
	Token  string
	UserID string
}
