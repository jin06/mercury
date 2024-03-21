package mqtt

import "io"

type Connect struct {
	Version         ProtocolVersion
	Clean           bool // v4
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

func (c *Connect) String() string {
	return ""
}

func (c *Connect) Decode(reader io.Reader) (err error) {
	buf := make([]byte, 16)
	if _, err = io.ReadFull(reader, buf); err != nil {
		return
	}
	c.Version = ProtocolVersion(buf[6])
	// var usernameFlag, passFlag, willRetain, willQoS, willFlag, cleanStart, reserved bool
	// usernameFlag := buf[8] & 0b10000000
	// passFlag := buf[8] & 0b01000000
	// willRetain := buf[8] & 0b00100000
	// willQoS := QoS(buf[8] & 0b00011000)
	// willFlag := buf[8] & 0b00000100
	c.Clean = buf[7]&0b00000010 == 0b00000010
	// reserved := buf[8] & 0b00000001
	c.KeepAlive = (uint16(buf[8]) << 8) + uint16(buf[9])
	return
}
