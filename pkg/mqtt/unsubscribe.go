package mqtt

import "fmt"

func NewUnsubscribe(header *FixedHeader, v ProtocolVersion) *Unsubscribe {
	return &Unsubscribe{
		BasePacket: &BasePacket{header, v},
		Properties: new(Properties),
	}
}

type Unsubscribe struct {
	*BasePacket
	PacketID     PacketID
	Properties   *Properties
	TopicFilters []string // dodo
}

func (u *Unsubscribe) String() string {
	return fmt.Sprintf("Unsubscribe - PacketID: %d, TopicFilters: %v", u.PacketID, u.TopicFilters)
}

func (u *Unsubscribe) Response() *Unsuback {
	resp := &Unsuback{
		BasePacket: newBasePacket(UNSUBACK, u.Version),
		PacketID:   u.PacketID,
		Properties: &Properties{},
	}
	resp.ReasonCodes = make([]ReasonCode, len(u.TopicFilters))
	return resp
}

func (u *Unsubscribe) Encode() ([]byte, error) {
	body, err := u.EncodeBody()
	if err != nil {
		return nil, err
	}
	u.FixedHeader.RemainingLength = VariableByteInteger(len(body))
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
	data, err := r.Read(u.Length())
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
	if u.Version == MQTT5 && u.Properties != nil {
		propertiesLength, _ := u.Properties.Encode()
		length += len(propertiesLength)
	}
	return length
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

	// Decode Properties (MQTT 5.0 only)
	if u.Version == MQTT5 {
		u.Properties = new(Properties)
		n, err := u.Properties.Decode(data[start:])
		if err != nil {
			return start, err
		}
		start += n
	}
	for start < len(data) {
		topic, n, err := decodeUTF8Str(data[start:])
		if err != nil {
			return start, err
		}
		u.TopicFilters = append(u.TopicFilters, topic)
		start += n
	}

	return len(data), nil
}

func (u *Unsubscribe) EncodeBody() ([]byte, error) {
	var data []byte

	// Encode Packet ID
	data = append(data, u.PacketID.Encode()...)

	// Encode Properties (MQTT 5.0 only)
	if u.Version == MQTT5 && u.Properties != nil {
		propertiesData, err := u.Properties.Encode()
		if err != nil {
			return nil, err
		}
		data = append(data, propertiesData...)
	}
	for _, topic := range u.TopicFilters {
		bytes, err := encodeUTF8Str(topic)
		if err != nil {
			return nil, err
		}
		data = append(data, bytes...)
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
