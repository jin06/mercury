package mqtt

import (
	"fmt"
)

const (
	CONNECT PacketType = iota + 1
	CONNACK
	PUBLISH
	PUBACK
	PUBREC
	PUBREL
	PUBCOMP
	SUBSCRIBE
	SUBACK
	UNSUBSCRIBE
	UNSUBACK
	PINGREQ
	PINGRESP
	DISCONNECT
	AUTH
)

type PacketType byte

func (pt *PacketType) Type() (string, error) {
	switch *pt {
	case CONNECT:
		return "CONNECT", nil
	case CONNACK:
		return "CONNACK", nil
	case PUBLISH:
		return "PUBLISH", nil
	case PUBACK:
		return "PUBACK", nil
	case PUBREC:
		return "PUBREC", nil
	case PUBCOMP:
		return "PUBCOMP", nil
	case SUBSCRIBE:
		return "SUBSCRIBE", nil
	case SUBACK:
		return "SUBACK", nil
	case UNSUBSCRIBE:
		return "UNSUBSCRIBE", nil
	case UNSUBACK:
		return "UNSUBACK", nil
	case PINGREQ:
		return "PINGREQ", nil
	case PINGRESP:
		return "PINGRESP", nil
	case DISCONNECT:
		return "DISCONNECT", nil
	case AUTH:
		return "AUTH", nil
	}
	return "", fmt.Errorf("unknown packet type %b", pt)
}

func toHeader(p PacketType) (result []byte) {
	return []byte{
		byte(p) << 4,
		0,
	}
}
