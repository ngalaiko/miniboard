package email

import (
	"fmt"
	"net/smtp"
	"net/url"

	"github.com/sirupsen/logrus"
)

// Config contains email client configuration.
type Config struct {
	Enabled  bool
	Addr     string
	From     string
	Username string
	Password string
}

// Client is used to send emails.
type Client struct {
	auth smtp.Auth
	cfg  *Config
}

// New returns new email client.
func New(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	if !cfg.Enabled {
		return &Client{cfg: cfg}, nil
	}

	addr, err := url.Parse(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	log().Infof("using %s smtp client", addr)

	return &Client{
		auth: smtp.PlainAuth(
			"",
			cfg.Username,
			cfg.Password,
			addr.Hostname(),
		),
		cfg: cfg,
	}, nil
}

// Send sends _payload_ to the client.
func (c *Client) Send(to string, subject string, payload string) error {
	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
		to, subject, payload))

	if !c.cfg.Enabled {
		log().Warnf("outgoing email: %s", msg)
		return nil
	}

	if err := smtp.SendMail(
		c.cfg.Addr,
		c.auth,
		c.cfg.From,
		[]string{to},
		msg,
	); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "email",
	})
}
