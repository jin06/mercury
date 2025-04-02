package mqtt

import "fmt"

func NewPubrel(header *FixedHeader, v ProtocolVersion) *Pubrel {
	return &Pubrel{
		BasePacket: &BasePacket{header, v},
		Properties: new(Properties),
	}
}

type Pubrel struct {
	*BasePacket
	PacketID PacketID
	ReasonCode
	Properties *Properties
}

func (p *Pubrel) String() string {
	return fmt.Sprintf("Pubrel - PacketID: %v, ReasonCode: %v", p.PacketID, p.ReasonCode)
}

func (p *Pubrel) ID() PacketID {
	return p.PacketID
}

func (p *Pubrel) Response() (resp Packet) {
	resp = &Pubcomp{
		BasePacket: newBasePacket(PUBCOMP, p.Version),
		PacketID:   p.PacketID,
		ReasonCode: V5_SUCCESS,
	}
	return
}

func (p *Pubrel) Encode() ([]byte, error) {
	body, err := p.EncodeBody()
	if err != nil {
		return nil, err
	}
	p.FixedHeader.Flags = 0b0010
	p.FixedHeader.RemainingLength = VariableByteInteger(len(body))
	header, err := p.FixedHeader.Encode()
	if err != nil {
		return nil, err
	}
	return append(header, body...), nil
}

func (p *Pubrel) Decode(data []byte) (int, error) {
	n, err := p.FixedHeader.Decode(data)
	if err != nil {
		return 0, err
	}
	bodyLen, err := p.DecodeBody(data[n:])
	return bodyLen + n, err
}

func (p *Pubrel) DecodeBody(data []byte) (int, error) {
	var start int

	// Decode Packet ID
	if err := p.PacketID.Decode(data); err != nil {
		return start, err
	} else {
		start += 2
	}

	// Decode Properties (MQTT 5.0 only)
	if p.Version == MQTT5 {
		p.ReasonCode = ReasonCode(data[start])
		start++
		if len(data) > start {
			n, err := p.Properties.Decode(data[start:])
			if err != nil {
				return start, err
			}
			start += n
		}
	}

	return len(data), nil
}

func (p *Pubrel) ReadBody(r *Reader) error {
	data, err := r.Read(p.Length())
	if err != nil {
		return err
	}
	_, err = p.DecodeBody(data)
	return err
}

func (p *Pubrel) Write(w *Writer) error {
	data, err := p.Encode()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (p *Pubrel) EncodeBody() ([]byte, error) {
	var data []byte

	// Encode Packet ID
	data = append(data, encodePacketID(p.PacketID)...)

	if p.Version == MQTT5 {
		data = append(data, byte(p.ReasonCode))
		propertiesData, err := p.Properties.Encode()
		if err != nil {
			return nil, err
		}
		data = append(data, propertiesData...)
	}

	return data, nil
}

func (p *Pubrel) WriteBody(w *Writer) error {
	data, err := p.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (p *Pubrel) PacketType() PacketType {
	return PUBREL
}

func (p *Pubrel) RemainingLength() int {
	length := 2 // Packet ID length
	if p.Version == MQTT5 && p.Properties != nil {
		propertiesLength, _ := p.Properties.Encode()
		length += len(propertiesLength)
	}
	return length
}
