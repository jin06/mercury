package mqtt

import "crypto/x509"

type Connect struct {
	Version         ProtocolVersion
	CleanSession    bool // v4
	KeepAlive       uint16
	ClientID        string
	Username        string
	Password        string
	WillQoS         bool
	WillFlag        bool
	WillRetain      bool // v4
	WillDelay       uint32
	WillProperties  map[string]string
	WillTopic       string
	WillMessage     []byte
	MessageExpiry   uint
	PayloadFormat   bool
	ContentType     string
	ResponseTopic   string
	CorrelationData []byte
	UserProperty    map[string]string
}
