package clients

type Server interface {
	On(c *Client)
}
