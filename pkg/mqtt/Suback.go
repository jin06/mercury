package mqtt

func NewSuback(header *FixedHeader, v ProtocolVersion) *Suback {
	return &Suback{BasePacket: &BasePacket{header, v}}
}

type Suback struct {
	*BasePacket
	PacketID    PacketID
	Properties  *Properties
	ReasonCodes []ReasonCode
}

func (s *Suback) Encode() ([]byte, error) {
	body, err := s.EncodeBody()
	if err != nil {
		return nil, err
	}
	s.FixedHeader.RemainingLength = len(body)
	header, err := s.FixedHeader.Encode()
	if err != nil {
		return nil, err
	}
	return append(header, body...), nil
}

func (s *Suback) EncodeBody() ([]byte, error) {
	var data []byte

	data = append(data, s.PacketID.Encode()...)
	// Encode Properties (MQTT 5.0 only)
	if s.Version == MQTT5 {
		propertiesData, err := s.Properties.Encode()
		if err != nil {
			return nil, err
		}
		data = append(data, propertiesData...)
	}

	for _, code := range s.ReasonCodes {
		data = append(data, byte(code))
	}

	return data, nil
}

func (s *Suback) Decode(data []byte) (int, error) {
	n, err := s.FixedHeader.Decode(data)
	if err != nil {
		return 0, err
	}
	bodyLen, err := s.DecodeBody(data[n:])
	return bodyLen + n, err
}

func (s *Suback) DecodeBody(data []byte) (int, error) {

	// Decode Properties (MQTT 5.0 only)
	s.Properties = new(Properties)
	start, err := s.Properties.Decode(data)
	if err != nil {
		return start, err
	}
	for len(data) > start {
		reason := ReasonCode(data[start])
		s.ReasonCodes = append(s.ReasonCodes, reason)
		start++
	}
	return len(data), nil
}

func (s *Suback) ReadBody(r *Reader) error {
	data, err := r.Read(s.FixedHeader.RemainingLength)
	if err != nil {
		return err
	}
	_, err = s.DecodeBody(data)
	return err
}

func (s *Suback) Write(w *Writer) error {
	data, err := s.Encode()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (s *Suback) WriteBody(w *Writer) error {
	data, err := s.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (s *Suback) PacketType() PacketType {
	return SUBACK
}

func (s *Suback) RemainingLength() int {
	length := 3 // Packet ID and Granted QoS length
	if s.Properties != nil {
		propertiesLength, _ := s.Properties.Encode()
		length += len(propertiesLength)
	}
	length += len(s.ReasonCodes)
	return length
}

func (s *Suback) String() string {
	return "Suback Packet"
}
