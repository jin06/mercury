package model

import (
	"time"

	"github.com/jin06/mercury/pkg/mqtt"
)

type Session struct {
	ClientID    string     `json:"client_id"`
	Will        *mqtt.Will `json:"will"`
	ConnectTime time.Time  `json:"connect_time"`
	//KeepTime last keep time or messaging time
	KeepTime time.Time `json:"keep_time"`
	Username string    `json:"username"`
	Clean    bool      `json:"clean"`
	// Session Expiry Interval inseconds
	Expiry uint32 `json:"expiry"`
}
