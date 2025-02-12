package mqtt

func NewUnsubscribe(header *FixedHeader) *Unsubscribe {
	return &Unsubscribe{FixedHeader: header}
}

type Unsubscribe struct {
	*FixedHeader
	QoS        QoS
	Dup        bool
	PacketID   PacketID
	Payload    []TopicWildcard
	Properties *Properties
}

func (u *Unsubscribe) Encode() ([]byte, error) {
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

func (u *Unsubscribe) Decode(data []byte) (int, error) {
	n, err := u.FixedHeader.Decode(data)
	if err != nil {
		return 0, err
	}
	bodyLen, err := u.DecodeBody(data[n:])
	return bodyLen + n, err
}

func (u *Unsubscribe) ReadBody(r *Reader) error {
	data, err := r.Read(u.FixedHeader.RemainingLength)
	if err != nil {
		return err
	}
	_, err = u.DecodeBody(data)
	return err
}

func (u *Unsubscribe) Write(w *Writer) error {
	data, err := u.Encode()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (u *Unsubscribe) PacketType() PacketType {
	return UNSUBSCRIBE
}

func (u *Unsubscribe) RemainingLength() int {
	length := 2 // Packet ID length
	for _, topic := range u.Payload {
		length += 2 + len(topic) // UTF-8 encoded topic length
	}
	return length
}

func (u *Unsubscribe) String() string {
	return "Unsubscribe Packet"
}

func (u *Unsubscribe) DecodeBody(data []byte) (int, error) {
	var start int

	// Decode Packet ID
	packetID, err := decodeUint16(data[start : start+2])
	if err != nil {
		return start, err
	}
	u.PacketID = PacketID(packetID)
	start += 2

	// Decode Payload
	for start < len(data) {
		topicWildcard, n, err := decodeUTF8Str(data[start:])
		if err != nil {
			return start, err
		}
		u.Payload = append(u.Payload, TopicWildcard(topicWildcard))
		start += n
	}

	return len(data), nil
}

func (u *Unsubscribe) EncodeBody() ([]byte, error) {
	var data []byte

	// Encode Packet ID
	data = append(data, encodePacketID(u.PacketID)...)

	// Encode Payload
	for _, topic := range u.Payload {
		if topicData, err := encodeUTF8Str(string(topic)); err != nil {
			return nil, err
		} else {
			data = append(data, topicData...)
		}
	}

	return data, nil
}

func (u *Unsubscribe) WriteBody(w *Writer) error {
	data, err := u.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
