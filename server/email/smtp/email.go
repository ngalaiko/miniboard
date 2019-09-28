package smtp

import (
	"fmt"
	"net/smtp"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Client is an email client.
type Client struct {
	auth smtp.Auth
	host string
	from string
}

// New returns new smtp client.
func New(host, port, from, username, password string) *Client {
	log("email").Infof("using %s:%s smtp client", host, port)
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

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
