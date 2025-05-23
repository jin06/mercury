package mqtt

import "fmt"

func NewSubscribe(header *FixedHeader, v ProtocolVersion) *Subscribe {
	return &Subscribe{
		BasePacket: &BasePacket{header, v},
		Properties: new(Properties),
	}
}

type Subscribe struct {
	*BasePacket
	PacketID      PacketID
	Subscriptions []*Subscription
	Properties    *Properties
}

func (s *Subscribe) String() string {
	subscriptions := make([]string, len(s.Subscriptions))
	for i, sub := range s.Subscriptions {
		subscriptions[i] = fmt.Sprintf("{TopicFilter: %s, QoS: %d}", sub.TopicFilter, sub.QoS)
	}
	return fmt.Sprintf("Subscribe - PacketID: %d, Subscriptions: [%s]", s.PacketID, fmt.Sprint(subscriptions))
}

func (s *Subscribe) ID() PacketID {
	return s.PacketID
}

func (s *Subscribe) Response() *Suback {
	resp := &Suback{
		BasePacket:  newBasePacket(SUBACK, s.Version),
		ReasonCodes: make([]ReasonCode, len(s.Subscriptions)),
		Properties:  &Properties{},
		PacketID:    s.PacketID,
	}
	return resp
}

func (s *Subscribe) Encode() ([]byte, error) {
	body, err := s.EncodeBody()
	if err != nil {
		return nil, err
	}

	s.FixedHeader.Flags = 0b0010

	s.FixedHeader.RemainingLength = VariableByteInteger(len(body))
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
	data, err := r.Read(s.Length())
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

func (s *Subscribe) DecodeBody(data []byte) (int, error) {
	var start int
	id, err := decodeUint16(data[:2])
	if err != nil {
		return start, err
	}
	start += 2
	s.PacketID = PacketID(id)
	// Decode Properties (MQTT 5.0 only)
	if s.Version == MQTT5 {
		s.Properties = new(Properties)
		n, err := s.Properties.Decode(data[start:])
		if err != nil {
			return start, err
		}
		start += n
	}

	// Decode Payload
	for start < len(data) {
		subscription := &Subscription{}
		if n, err := subscription.Decode(data[start:]); err != nil {
			return start, err
		} else {
			start += n
		}
		s.Subscriptions = append(s.Subscriptions, subscription)
	}

	return len(data), nil
}

func (s *Subscribe) EncodeBody() ([]byte, error) {
	var data []byte

	data = encodeUint16(uint16(s.PacketID))

	// Encode Properties (MQTT 5.0 only)
	if s.Version == MQTT5 && s.Properties != nil {
		propertiesData, err := s.Properties.Encode()
		if err != nil {
			return nil, err
		}
		data = append(data, propertiesData...)
	}

	// Encode Payload
	for _, subscription := range s.Subscriptions {
		if topicData, err := encodeUTF8Str(string(subscription.TopicFilter)); err != nil {
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
	TopicFilter       string
	RetainHandling    byte
	RetainAsPublished bool
	NoLocal           bool
	QoS               QoS
}

func (s *Subscription) Encode() []byte {
	return nil
}

func (s *Subscription) Decode(data []byte) (int, error) {
	topic, n, err := decodeUTF8Str(data)
	if err != nil {
		return n, err
	}
	if n >= len(data) {
		return n, ErrMalformedPacket
	}
	s.TopicFilter = topic
	options := data[n]
	s.QoS = QoS(options & 0b00000011)
	s.NoLocal = (options & 0b00000100) != 0
	s.RetainAsPublished = (options & 0b00001000) != 0
	n++
	return n, nil
}
