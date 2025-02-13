package clients

type Server interface {
	Reg(c *Client) error
}
