package mqtt

type Publish struct {
	*FixedHeader
	PacketID               PacketID
	Dup                    bool
	Qos                    QoS
	Retain                 bool
	Topic                  string
	Payload                []byte
	PayloadFormatIndicator bool   // from mqtt5
	MessageExpiryInterval  uint16 // from mqtt5 (seconds)
	TopicAlias             uint16 // from mqtt5
	ResponseTopic          string // from mqtt5
	CorrelationData        []byte // from mqtt5
	Properties             *Properties
	SubscriptionIdentifier []uint32
	ContentType            string
	Version                ProtocolVersion
}

func NewPublish(header *FixedHeader) *Publish {
	return &Publish{FixedHeader: header}
}

func (p *Publish) Encode() ([]byte, error) {
	// write header
	result := toHeader(PUBLISH)
	if p.Dup {
		result[0] |= 0b00001000
	}
	result[0] |= (byte(p.Qos) << 1)
	// write topic name
	if bytes, err := encodeUTF8Str(p.Topic); err != nil {
		return nil, err
	} else {
		result = append(result, bytes...)
	}
	// write message ID
	result = append(result, p.PacketID.Encode()...)

	if p.Version == MQTT5 {
		if p.Properties != nil {
			bytes, err := p.Properties.Encode()
			if err != nil {
				return nil, err
			}
			result = append(result, bytes...)
		}
	}
	result = append(result, p.Payload...)

	return result, nil
}

func (p *Publish) DecodeBody(data []byte) (int, error) {
	var start int

	// Decode Topic
	topic, n, err := decodeUTF8Str(data[start:])
	if err != nil {
		return start, err
	}
	p.Topic = topic
	start += n

	// Decode Packet ID
	packetID, err := decodeUint16(data[start : start+2])
	if err != nil {
		return start, err
	}
	p.PacketID = PacketID(packetID)
	start += 2

	// Decode Properties (MQTT 5.0 only)
	if p.Version == MQTT5 {
		properties := &Properties{}
		n, err := properties.Decode(data[start:])
		if err != nil {
			return start, err
		}
		p.Properties = properties
		start += n
	}

	// Decode Payload
	p.Payload = data[start:]

	return len(data), nil
}

func (p *Publish) ReadBody(r *Reader) error {
	data, err := r.Read(p.FixedHeader.RemainingLength)
	if err != nil {
		return err
	}
	_, err = p.DecodeBody(data)
	return err
}

func (p *Publish) Write(w *Writer) error {
	data, err := p.Encode()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (p *Publish) EncodeBody() ([]byte, error) {
	var data []byte

	// Encode Topic
	if topicData, err := encodeUTF8Str(p.Topic); err != nil {
		return nil, err
	} else {
		data = append(data, topicData...)
	}

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

	// Encode Payload
	data = append(data, p.Payload...)

	return data, nil
}

func (p *Publish) WriteBody(w *Writer) error {
	data, err := p.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
