package mqtt

func NewPuback(header *FixedHeader, v ProtocolVersion) *Puback {
	return &Puback{BasePacket: &BasePacket{header, v}}
}

type Puback struct {
	*BasePacket
	PacketID     PacketID
	ReasonCode   ReasonCode
	ReasonString string
	Properties   *Properties
}

func (p *Puback) Encode() ([]byte, error) {
	result := toHeader(PUBACK)
	if p.Version == MQTT5 {
		result[1] = 0b00000100
		result = append(result, p.PacketID.Encode()...)
		result = append(result, byte(p.ReasonCode), 0)
		var length byte = 0
		if reasonString, err := encodeUTF8Str(p.ReasonString); err != nil {
			return []byte{}, nil
		} else {
			length += byte(len(reasonString))
			result = append(result, reasonString...)
		}

		if userProperties, err := p.Properties.Encode(); err != nil {
			return []byte{}, nil
		} else {
			length += byte(len(userProperties))
			result = append(result, userProperties...)
		}
		result[5] = length

	} else {
		result[1] = 0b01000000
		result = append(result, p.PacketID.Encode()...)
	}

	return result, nil
}

func (p *Puback) Decode(data []byte) (int, error) {
	n, err := p.FixedHeader.Decode(data)
	if err != nil {
		return 0, err
	}
	bodyLen, err := p.DecodeBody(data[n:])
	return bodyLen + n, err
}

func (p *Puback) DecodeBody(data []byte) (int, error) {
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

func (p *Puback) ReadBody(r *Reader) error {
	data, err := r.Read(p.Length())
	if err != nil {
		return err
	}
	_, err = p.DecodeBody(data)
	return err
}

func (p *Puback) Write(w *Writer) error {
	data, err := p.Encode()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (p *Puback) EncodeBody() ([]byte, error) {
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

func (p *Puback) WriteBody(w *Writer) error {
	data, err := p.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (p *Puback) PacketType() PacketType {
	return PUBACK
}

func (p *Puback) RemainingLength() int {
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

func (p *Puback) String() string {
	return "Puback Packet"
}
