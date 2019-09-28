package smtp

import (
	"fmt"
	"net/smtp"

	"github.com/pkg/errors"
)

// Client is an email client.
type Client struct {
	auth smtp.Auth
	host string
	from string
}

// New returns new smtp client.
func New(host, port, from, username, password string) *Client {
	return &Client{
		auth: smtp.PlainAuth(
			"",
			username,
			password,
			host,
		),
		host: fmt.Sprintf("%s:%s", host, port),
		from: from,
	}
}

// Send sends _payload_ to the client.
func (c *Client) Send(to string, subject string, payload string) error {
	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
		to, subject, payload))
	if err := smtp.SendMail(
		c.host,
		c.auth,
		c.from,
		[]string{to},
		msg,
	); err != nil {
		return errors.Wrap(err, "failed to send email")
	}
	return nil
}
