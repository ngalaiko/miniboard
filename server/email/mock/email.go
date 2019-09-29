package mock

// Client is an email client.
type Client struct {
	msg []string
}

// New returns new smtp client.
func New() *Client {
	return &Client{}
}

// Send sends _payload_ to the client.
func (c *Client) Send(to string, subject string, payload string) error {
	c.msg = append(c.msg, payload)
	return nil
}

// LastMessage returns the last received message.
func (c *Client) LastMessage() string {
	if len(c.msg) == 0 {
		return ""
	}
	return c.msg[len(c.msg)-1]
}
