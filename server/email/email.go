package email

// Client an email client.
type Client interface {
	// Send sends an email message.
	Send(to string, subject string, payload string) error
}
