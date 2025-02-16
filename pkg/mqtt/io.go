package mqtt

import (
	"fmt"
)

// Read FixedHeader and remaining length
func ReadPacket(reader *Reader) (packet Packet, err error) {
	header := &FixedHeader{}
	if err := header.Read(reader); err != nil {
		return nil, err
	}
	switch header.PacketType {
	case CONNECT:
		packet = NewConnect(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	case CONNACK:
		packet = NewConnack(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	case PUBLISH:
		packet = NewPublish(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	case PUBACK:
		packet = NewPuback(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	case PUBREC:
		packet = NewPubrec(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	case PUBREL:
		packet = NewPubrel(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	case PUBCOMP:
		packet = NewPubcomp(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	case SUBSCRIBE:
		packet = NewSubscribe(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	case SUBACK:
		packet = NewSuback(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	case UNSUBSCRIBE:
		packet = NewUnsubscribe(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	case UNSUBACK:
		packet = NewUnsuback(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	case PINGREQ:
		packet = NewPingreq(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	case PINGRESP:
		packet = NewPingresp(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	case DISCONNECT:
		packet = NewDisconnect(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	case AUTH:
		packet = NewAuth(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	default:
		fmt.Println(header)
	}
	return
}

func WritePacket(writer *Writer, p Packet) error {
	return p.Write(writer)
}
