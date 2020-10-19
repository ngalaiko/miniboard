package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

// Config contains email client configuration.
type Config struct {
	Enabled  bool
	Addr     string
	From     string
	Username string
	Password string
}

type logger interface {
	Info(string, ...interface{})
	Warn(string, ...interface{})
}

// Client is used to send emails.
type Client struct {
	logger logger
	auth   smtp.Auth
	cfg    *Config
}

// New returns new email client.
func New(cfg *Config, logger logger) (*Client, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	if !cfg.Enabled {
		return &Client{
			logger: logger,
			cfg:    cfg,
		}, nil
	}

	return &Client{
		logger: logger,
		auth: smtp.PlainAuth(
			"",
			cfg.Username,
			cfg.Password,
			strings.Split(cfg.Addr, ":")[0],
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
		c.logger.Warn("outgoing email: %s", msg)
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
