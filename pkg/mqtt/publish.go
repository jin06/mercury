package mqtt

import (
	"fmt"
)

func NewPublish(header *FixedHeader, v ProtocolVersion) *Publish {
	p := &Publish{
		BasePacket: &BasePacket{header, v},
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

func (p *Publish) Clone() *Publish {
	clone := &Publish{
		BasePacket: p.BasePacket.Clone(),
		PacketID:   p.PacketID,
		Dup:        p.Dup,
		Qos:        p.Qos,
		Retain:     p.Retain,
		Topic:      p.Topic,
		Payload:    append([]byte{}, p.Payload...), // Deep copy of Payload
		Properties: p.Properties.Clone(),
	}
	return clone
}

func (p *Publish) String() string {
	return fmt.Sprintf("Publish - Dup: %t, Qos: %d, Retain: %t, Topic: %s, PacketID: %d, Payload: %s",
		p.Dup, p.Qos, p.Retain, p.Topic, p.PacketID, p.Payload)
}

func (p *Publish) ID() PacketID {
	return p.PacketID
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
	b := 0b110 & (byte(p.Qos) << 1)
	if p.Retain {
		b |= 0b1
	}
	if p.Dup {
		b |= 0b1000
	}
	return b
}

func (p *Publish) fromFlags() {
	p.Dup = p.FixedHeader.Flags&0b1000 != 0
	p.Qos = QoS((p.FixedHeader.Flags & 0b110) >> 1)
	p.Retain = p.FixedHeader.Flags&0b1 != 0
}

func (p *Publish) Encode() ([]byte, error) {
	body, err := p.EncodeBody()
	if err != nil {
		return nil, err
	}
	p.FixedHeader.Flags = p.flags()
	p.RemainingLength = VariableByteInteger(len(body))

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

	if !p.Qos.Zero() {
		if err := p.PacketID.Decode(data[start:]); err != nil {
			return start, err
		} else {
			start += 2
		}
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
	data := make([]byte, 0)

	// Encode Topic
	if topicData, err := p.Topic.Encode(); err != nil {
		return nil, err
	} else {
		data = append(data, topicData...)
	}

	// Encode Packet ID
	if !p.Qos.Zero() {
		data = append(data, p.PacketID.Encode()...)
	}
	// Encode Properties (MQTT 5.0 only)
	if p.Version == MQTT5 && p.Properties != nil {
		if propertiesData, err := p.Properties.Encode(); err != nil {
			return nil, err
		} else {
			data = append(data, propertiesData...)
		}
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
