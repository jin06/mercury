package mqtt

import "fmt"

func NewPubcomp(header *FixedHeader, v ProtocolVersion) *Pubcomp {
	return &Pubcomp{
		BasePacket: &BasePacket{header, v},
		Properties: new(Properties),
	}
}

type Pubcomp struct {
	*BasePacket
	ReasonCode ReasonCode
	Properties *Properties
	PacketID   PacketID
}

func (p *Pubcomp) PID() PacketID {
	return p.PacketID
}

func (p *Pubcomp) String() string {
	return fmt.Sprintf("Pubcomp - PacketID: %v, ReasonCode: %v, Properties: %v", p.PacketID, p.ReasonCode, p.Properties)
}

func (p *Pubcomp) ID() PacketID {
	return p.PacketID
}

func (p *Pubcomp) Encode() ([]byte, error) {
	body, err := p.EncodeBody()
	if err != nil {
		return nil, err
	}
	p.FixedHeader.RemainingLength = VariableByteInteger(len(body))
	header, err := p.FixedHeader.Encode()
	if err != nil {
		return nil, err
	}
	return append(header, body...), nil
}

func (p *Pubcomp) Decode(data []byte) (int, error) {
	n, err := p.FixedHeader.Decode(data)
	if err != nil {
		return 0, err
	}
	bodyLen, err := p.DecodeBody(data[n:])
	return bodyLen + n, err
}

func (p *Pubcomp) DecodeBody(data []byte) (int, error) {
	var start int

	// Decode Packet ID
	packetID, err := decodeUint16(data[start : start+2])
	if err != nil {
		return start, err
	}
	p.PacketID = PacketID(packetID)
	start += 2

	// Decode Reason Code (MQTT 5.0 only)
	if p.Version == MQTT5 {
		if len(data) > start {
			p.ReasonCode = ReasonCode(data[start])
			start++
		}

		// Decode Properties (MQTT 5.0 only)
		if len(data) > start {
			p.Properties = new(Properties)
			n, err := p.Properties.Decode(data[start:])
			if err != nil {
				return start, err
			}
			start += n
		}
	}

	return len(data), nil
}

func (p *Pubcomp) ReadBody(r *Reader) error {
	data, err := r.Read(p.Length())
	if err != nil {
		return err
	}
	_, err = p.DecodeBody(data)
	return err
}

func (p *Pubcomp) Write(w *Writer) error {
	data, err := p.Encode()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (p *Pubcomp) EncodeBody() ([]byte, error) {
	var data []byte

	// Encode Packet ID
	data = append(data, encodePacketID(p.PacketID)...)

	// Encode Reason Code (MQTT 5.0 only)
	if p.Version == MQTT5 {
		data = append(data, byte(p.ReasonCode))

		// Encode Properties (MQTT 5.0 only)
		if p.Properties != nil {
			propertiesData, err := p.Properties.Encode()
			if err != nil {
				return nil, err
			}
			data = append(data, propertiesData...)
		}
	}

	return data, nil
}

func (p *Pubcomp) WriteBody(w *Writer) error {
	data, err := p.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (p *Pubcomp) PacketType() PacketType {
	return PUBCOMP
}

func (p *Pubcomp) RemainingLength() int {
	length := 2 // Packet ID length
	if p.Version == MQTT5 {
		length++ // Reason Code length
		if p.Properties != nil {
			propertiesLength, _ := p.Properties.Encode()
			length += len(propertiesLength)
		}
	}
	return length
}
