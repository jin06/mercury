package mqtt

import "io"

type Disconnect struct {
	Version    ProtocolVersion
	ResionCode ReasonCode
}

func (d *Disconnect) Encode() (result []byte, err error) {
	result = toHeader(DISCONNECT)
	if d.Version == MQTT5 {
		result = append(result, byte(d.ResionCode))
	}
	return
}

func (d *Disconnect) Decode(reader io.Reader) (err error) {
	if d.Version == MQTT5 {
		buf := make([]byte, 1)
		if _, err = io.ReadFull(reader, buf); err != nil {
			return
		}
		d.ResionCode = ReasonCode(buf[0])
	}
	return nil
}
