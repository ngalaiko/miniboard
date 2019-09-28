package disabled

import (
	"errors"

	"github.com/sirupsen/logrus"
)

// Client is an email client.
type Client struct {
}

// New returns new smtp client.
func New() *Client {
	log("email").Warning("disabled")
	return &Client{}
}

// Send sends _payload_ to the client.
func (c *Client) Send(to string, subject string, payload string) error {
	return errors.New("email client disabled")
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
