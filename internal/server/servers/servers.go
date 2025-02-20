package servers

import "github.com/jin06/mercury/internal/server"

func NewServer() server.Server {
	return newGeneric()
}
