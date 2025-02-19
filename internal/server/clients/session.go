package clients

import (
	"time"

	"github.com/jin06/mercury/pkg/mqtt"
)

type Session struct {
	ClinetID      string `json:"client_id"`
	ConnectedTime time.Time
	MQTTVersion   mqtt.ProtocolVersion
}
