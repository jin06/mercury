package model

import "time"

type Session struct {
	ClientID    string    `json:"client_id"`
	Will        Message   `json:"will"`
	ConnectTime time.Time `json:"connect_time"`
	//KeepTime last keep time or messaging time
	KeepTime time.Time `json:"keep_time"`
	Username string    `json:"username"`
	Clean    bool      `json:"clean"`
	// Session Expiry Interval inseconds
	Expiry uint32 `json:"expiry"`
}
