package mqtt

import (
	"errors"
	"fmt"
	"io"
)

type Connect struct {
	Version                    ProtocolVersion
	Clean                      bool // v4
	KeepAlive                  uint16
	ClientID                   []byte
	Username                   string
	Password                   string
	WillQoS                    bool
	WillFlag                   bool
	WillRetain                 bool // v4
	WillDelay                  uint32
	WillProperties             map[string]string
	WillTopic                  []byte
	WillMessage                []byte
	willDelayInterval          uint64 // only mqtt5
	MessageExpiry              uint64
	SessionExpiryInterval      uint64 // mqtt5
	RequestResponseInformation byte
	RequestProblemInformation  byte
	PayloadFormat              bool
	ContentType                string
	ResponseTopic              string
	CorrelationData            []byte
	UserProperty               map[string]string
}

func (c *Connect) String() string {
	return ""
}

func (c *Connect) Decode(reader io.Reader) (err error) {
	// buf := make([]byte, 16)

	protocolName, err := decodeProtocolName(reader)
	if err != nil {
		return
	}
	fmt.Println(protocolName)
	buf := make([]byte, 4)
	if _, err = reader.Read(buf); err != nil {
		return
	}

	c.Version = ProtocolVersion(buf[0])
	// var usernameFlag, passFlag, willRetain, willQoS, willFlag, cleanStart, reserved bool
	usernameFlag := buf[1]&0b10000000 == 0b10000000
	passFlag := buf[1]&0b01000000 == 0b01000000
	// willRetain := buf[8] & 0b00100000
	// willQoS := QoS(buf[8] & 0b00011000)
	willFlag := buf[1]&0b00000100 == 0b00000100
	c.Clean = buf[1]&0b00000010 == 0b00000010
	// reserved := buf[8] & 0b00000001
	// measured in seconds
	c.KeepAlive = decodeKeepAlive(buf[2:])
	c.decodeProperties(reader)
	if c.ClientID, err = decodeUTF8(reader); err != nil {
		return
	}
	switch c.Version {
	case MQTT3, MQTT4:
		{
			if willFlag {
				if c.WillTopic, err = decodeUTF8(reader); err != nil {
					return
				}
				if c.WillMessage, err = decodeUTF8(reader); err != nil {
					return
				}
			}
			if usernameFlag {
				if c.Username, err = decodeUTF8Str(reader); err != nil {
					return
				}
			}
			if passFlag {
				if c.Password, err = decodeUTF8Str(reader); err != nil {
					return
				}

			}
		}
	case MQTT5:
		{
			// read property length
			proLen := make([]byte, 1)
			if _, err = reader.Read(proLen); err != nil {
				return
			}
			willProperties := make([]byte, int(proLen[0]))
			if _, err = reader.Read(willProperties); err != nil {
				return
			}
			for {
				if len(willProperties) <= 0 {
					break
				}
				idetifier := willProperties[0]
				switch idetifier {
				case 0x18:
					{
						if c.willDelayInterval, err = bytesToUint64(willProperties[1:5]); err != nil {
							return
						}
					}
					// payload format
				case 0x01:
					{

					}
					// publication expiry interval
				case 0x02:
					{

					}
					// content type
				case 0x03:
					{

					}
					// reponse topic
				case 0x08:
					{

					}
					// correlation data
				case 0x09:
					{

					}
				}

			}

			if c.willDelayInterval, err = bytesToUint64(willProperties[0:4]); err != nil {
				return
			}
		}
	}
	return
}

func (c *Connect) decodeProperties(reader io.Reader) (err error) {
	if c.Version == MQTT5 {
		prolength := make([]byte, 1)
		if _, err = reader.Read(prolength); err != nil {
			return
		}
		properties := make([]byte, prolength[0])
		if _, err = reader.Read(properties); err != nil {
			return
		}
		for {
			if len(properties) <= 0 {
				break
			}
			switch properties[0] {
			case 0x11:
				{
					if len(properties) < 5 {
						return errors.New("protocol error")
					}
					if c.SessionExpiryInterval, err = bytesToUint64(properties[1:5]); err != nil {
						return
					}
					properties = properties[5:]
				}
			case 0x19:
				{
					if len(properties) < 2 {
						return errors.New("protocol error")
					}

					c.RequestResponseInformation = properties[1]
					properties = properties[2:]
				}
			case 0x17:
				{
					if len(properties) < 2 {
						return errors.New("protocol error")
					}
					c.RequestProblemInformation = properties[1]
					properties = properties[2:]
				}
				// receive max
			case 0x21:
				{

				}
				// Max packet size
			case 0x27:
				{

				}
				//  Topic Alias Max
			case 0x22:
				{
				}
				// User properties
			case 0x26:
				{
					break
				}
			}
		}
	}
	return
}
