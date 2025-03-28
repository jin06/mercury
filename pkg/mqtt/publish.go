package mqtt

import (
	"fmt"
)

func NewPublish(header *FixedHeader, v ProtocolVersion) *Publish {
	p := &Publish{
		BasePacket: &BasePacket{
			FixedHeader: header,
			Version:     v},
		Properties: new(Properties),
	}
	p.fromFlags()
	return p
}

type Publish struct {
	*BasePacket
	PacketID   PacketID
	Dup        bool
	Qos        QoS
	Retain     bool
	Topic      UTF8String
	Payload    []byte
	Properties *Properties
}

func (p *Publish) String() string {
	return fmt.Sprintf("Publish - Dup: %t, Qos: %d, Retain: %t, Topic: %s, PacketID: %d, Payload: %s",
		p.Dup, p.Qos, p.Retain, p.Topic, p.PacketID, p.Payload)
}

func (p *Publish) Response() (resp Packet, err error) {
	switch p.Qos {
	case QoS0:
		resp = nil
	case QoS1:
		resp = &Puback{
			BasePacket: newBasePacket(PUBACK, p.Version),
			PacketID:   p.PacketID,
			ReasonCode: V5_SUCCESS,
		}
	case QoS2:
		resp = &Pubrec{
			BasePacket: newBasePacket(PUBREC, p.Version),
			PacketID:   p.PacketID,
			ReasonCode: V5_SUCCESS,
		}
	default:
		return nil, ErrInvalidQoS
	}
	return
}

func (p *Publish) flags() byte {
	b := 0b00000110 & byte(p.Qos) << 1
	if p.Retain {
		b |= 0b00000001
	}
	if p.Dup {
		b |= 0b00001000
	}
	return b
}

func (p *Publish) fromFlags() {
	p.Dup = p.FixedHeader.Flags&0b00001000 != 0
	p.Qos = QoS((p.FixedHeader.Flags & 0b00000110) >> 1)
	p.Retain = p.FixedHeader.Flags&0b00000001 != 0
}

func (p *Publish) Encode() ([]byte, error) {
	body, err := p.EncodeBody()
	if err != nil {
		return nil, err
	}
	p.RemainingLength = VariableByteInteger(len(body))
	p.FixedHeader.Flags = p.flags()

	header, err := p.FixedHeader.Encode()
	if err != nil {
		return nil, err
	}
	return append(header, body...), nil
}

func (p *Publish) DecodeBody(data []byte) (int, error) {
	var start int

	// Decode Topic
	if n, err := p.Topic.Decode(data[start:]); err != nil {
		return start, err
	} else {
		start += n
	}

	// Decode Packet ID
	if packetID, err := decodeUint16(data[start : start+2]); err != nil {
		return start, err
	} else {
		p.PacketID = PacketID(packetID)
		start += 2
	}

	// Decode Properties (MQTT 5.0 only)
	if p.Version == MQTT5 {
		if n, err := p.Properties.Decode(data[start:]); err != nil {
			return start, err
		} else {
			start += n
		}
	}

	// Decode Payload
	p.Payload = data[start:]

	return len(data), nil
}

func (p *Publish) ReadBody(r *Reader) error {
	data, err := r.Read(p.Length())
	if err != nil {
		return err
	}
	_, err = p.DecodeBody(data)
	return err
}

func (p *Publish) Write(w *Writer) error {
	data, err := p.Encode()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (p *Publish) EncodeBody() ([]byte, error) {
	var data []byte

	// Encode Topic
	if topicData, err := encodeUTF8Str(p.Topic); err != nil {
		return nil, err
	} else {
		data = append(data, topicData...)
	}

	// Encode Packet ID
	data = append(data, encodePacketID(p.PacketID)...)

	// Encode Properties (MQTT 5.0 only)
	if p.Version == MQTT5 && p.Properties != nil {
		propertiesData, err := p.Properties.Encode()
		if err != nil {
			return nil, err
		}
		data = append(data, propertiesData...)
	}

	// Encode Payload
	data = append(data, p.Payload...)

	return data, nil
}

func (p *Publish) WriteBody(w *Writer) error {
	data, err := p.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
