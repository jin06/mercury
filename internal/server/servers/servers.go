package servers

import (
	"github.com/jin06/mercury/internal/config"
	"github.com/jin06/mercury/internal/server"
)

func NewServer(mode config.Mode) server.Server {
	switch mode {
	case config.MemoryMode:
		return newGeneric()
	}
	return nil
}
