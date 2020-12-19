package jwt

// Token contains information about an authorized user.
type Token struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}
