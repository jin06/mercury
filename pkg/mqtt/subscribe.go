package mqtt

func NewSubscribe(header *FixedHeader) *Subscribe {
	return &Subscribe{FixedHeader: header}
}

type Subscribe struct {
	*FixedHeader
	Version       ProtocolVersion
	PacketID      PacketID
	TopicWildcard TopicWildcard
	Payload       []Subscription
	Properties    *Properties
}

func (s *Subscribe) Response() *Suback {
	return &Suback{}
}

func (s *Subscribe) Encode() ([]byte, error) {
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

func (s *Subscribe) Decode(data []byte) (int, error) {
	n, err := s.FixedHeader.Decode(data)
	if err != nil {
		return 0, err
	}
	bodyLen, err := s.DecodeBody(data[n:])
	return bodyLen + n, err
}

func (s *Subscribe) ReadBody(r *Reader) error {
	data, err := r.Read(s.FixedHeader.RemainingLength)
	if err != nil {
		return err
	}
	_, err = s.DecodeBody(data)
	return err
}

func (s *Subscribe) Write(w *Writer) error {
	data, err := s.Encode()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (s *Subscribe) PacketType() PacketType {
	return SUBSCRIBE
}

func (s *Subscribe) RemainingLength() int {
	length := 2 // Packet ID length
	if s.Version == MQTT5 && s.Properties != nil {
		propertiesLength, _ := s.Properties.Encode()
		length += len(propertiesLength)
	}
	length += len(s.TopicWildcard) + 2
	for _, subscription := range s.Payload {
		length += len(subscription.TopicWildcard) + 3
	}
	return length
}

func (s *Subscribe) String() string {
	return "Subscribe Packet"
}

func (s *Subscribe) DecodeBody(data []byte) (int, error) {
	var start int

	// Decode Packet ID
	packetID, err := decodeUint16(data[start : start+2])
	if err != nil {
		return start, err
	}
	s.PacketID = PacketID(packetID)
	start += 2

	// Decode Properties (MQTT 5.0 only)
	if s.Version == MQTT5 {
		s.Properties = new(Properties)
		n, err := s.Properties.Decode(data[start:])
		if err != nil {
			return start, err
		}
		start += n
	}

	// Decode Topic Wildcard
	topicWildcard, n, err := decodeUTF8Str(data[start:])
	if err != nil {
		return start, err
	}
	s.TopicWildcard = TopicWildcard(topicWildcard)
	start += n

	// Decode Payload
	for start < len(data) {
		subscription := Subscription{}
		topic, n, err := decodeUTF8Str(data[start:])
		if err != nil {
			return start, err
		}
		subscription.TopicWildcard = TopicWildcard(topic)
		start += n

		qos := QoS(data[start])
		subscription.QoS = qos
		start++

		s.Payload = append(s.Payload, subscription)
	}

	return len(data), nil
}

func (s *Subscribe) EncodeBody() ([]byte, error) {
	var data []byte

	// Encode Packet ID
	data = append(data, s.PacketID.Encode()...)

	// Encode Properties (MQTT 5.0 only)
	if s.Version == MQTT5 && s.Properties != nil {
		propertiesData, err := s.Properties.Encode()
		if err != nil {
			return nil, err
		}
		data = append(data, propertiesData...)
	}

	// Encode Topic Wildcard
	if topicData, err := encodeUTF8Str(string(s.TopicWildcard)); err != nil {
		return nil, err
	} else {
		data = append(data, topicData...)
	}

	// Encode Payload
	for _, subscription := range s.Payload {
		if topicData, err := encodeUTF8Str(string(subscription.TopicWildcard)); err != nil {
			return nil, err
		} else {
			data = append(data, topicData...)
		}
		data = append(data, byte(subscription.QoS))
	}

	return data, nil
}

func (s *Subscribe) WriteBody(w *Writer) error {
	data, err := s.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

type Subscription struct {
	TopicWildcard TopicWildcard
	QoS           QoS
}
