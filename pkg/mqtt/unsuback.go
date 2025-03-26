package mqtt

func NewUnsuback(header *FixedHeader, v ProtocolVersion) *Unsuback {
	return &Unsuback{BasePacket: &BasePacket{header, v}}
}

type Unsuback struct {
	*BasePacket
	PacketID    PacketID
	ReasonCodes []ReasonCode
	Properties  *Properties
}

func (u *Unsuback) Encode() ([]byte, error) {
	body, err := u.EncodeBody()
	if err != nil {
		return nil, err
	}
	u.FixedHeader.RemainingLength = len(body)
	header, err := u.FixedHeader.Encode()
	if err != nil {
		return nil, err
	}
	return append(header, body...), nil
}

func (u *Unsuback) Decode(data []byte) (int, error) {
	n, err := u.FixedHeader.Decode(data)
	if err != nil {
		return 0, err
	}
	bodyLen, err := u.DecodeBody(data[n:])
	return bodyLen + n, err
}

func (u *Unsuback) DecodeBody(data []byte) (int, error) {
	var start int

	// Decode Packet ID
	packetID, err := decodeUint16(data[start : start+2])
	if err != nil {
		return start, err
	}
	u.PacketID = PacketID(packetID)
	start += 2

	// Decode Reason Code (MQTT 5.0 only)
	if u.Version == MQTT5 {
		// Decode Properties (MQTT 5.0 only)
		if len(data) > start {
			u.Properties = new(Properties)
			n, err := u.Properties.Decode(data[start:])
			if err != nil {
				return start, err
			}
			start += n
		}
	}
	for len(data) > start {
		reason := ReasonCode(data[start])
		u.ReasonCodes = append(u.ReasonCodes, reason)
		start++
	}

	// Decode Payload
	return len(data), nil
}

func (u *Unsuback) ReadBody(r *Reader) error {
	data, err := r.Read(u.FixedHeader.RemainingLength)
	if err != nil {
		return err
	}
	_, err = u.DecodeBody(data)
	return err
}

func (u *Unsuback) Write(w *Writer) error {
	data, err := u.Encode()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (u *Unsuback) EncodeBody() ([]byte, error) {
	var data []byte

	// Encode Packet ID
	data = append(data, encodePacketID(u.PacketID)...)

	// Encode Reason Code (MQTT 5.0 only)
	if u.Version == MQTT5 {
		// Encode Properties (MQTT 5.0 only)
		if u.Properties != nil {
			propertiesData, err := u.Properties.Encode()
			if err != nil {
				return nil, err
			}
			data = append(data, propertiesData...)
		}
	}
	for _, code := range u.ReasonCodes {
		data = append(data, byte(code))
	}
	return data, nil
}

func (u *Unsuback) WriteBody(w *Writer) error {
	data, err := u.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (u *Unsuback) PacketType() PacketType {
	return UNSUBACK
}

func (u *Unsuback) RemainingLength() int {
	length := 2 // Packet ID length
	if u.Version == MQTT5 {
		length++ // Reason Code length
		if u.Properties != nil {
			propertiesLength, _ := u.Properties.Encode()
			length += len(propertiesLength)
		}
	}
	length += len(u.ReasonCodes)
	return length
}

func (u *Unsuback) String() string {
	return "Unsuback Packet"
}
