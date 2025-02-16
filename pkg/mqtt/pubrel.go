package mqtt

func NewPubrel(header *FixedHeader) *Pubrel {
	return &Pubrel{FixedHeader: header}
}

type Pubrel struct {
	*FixedHeader
	Version    ProtocolVersion
	PacketID   PacketID
	QoS        QoS
	Dup        bool
	Properties *Properties
}

func (p *Pubrel) Encode() ([]byte, error) {
	body, err := p.EncodeBody()
	if err != nil {
		return nil, err
	}
	p.FixedHeader.RemainingLength = len(body)
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
	packetID, err := decodeUint16(data[start : start+2])
	if err != nil {
		return start, err
	}
	p.PacketID = PacketID(packetID)
	start += 2

	// Decode Properties (MQTT 5.0 only)
	if p.Version == MQTT5 {
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

func (p *Pubrel) ReadBody(r *Reader) error {
	data, err := r.Read(p.FixedHeader.RemainingLength)
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

	// Encode Properties (MQTT 5.0 only)
	if p.Version == MQTT5 && p.Properties != nil {
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

func (p *Pubrel) String() string {
	return "Pubrecl Packet"
}
