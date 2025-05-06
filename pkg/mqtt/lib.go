package mqtt

import "regexp"

// ValidateMQTTTopic validates if the MQTT topic name is valid according to MQTT rules.
func ValidateMQTTTopic(topic string) bool {
	// Regular expression to validate MQTT topic format
	// This regex checks if the topic contains only alphanumeric characters, slashes, and special MQTT characters (+, #)
	re := regexp.MustCompile(`^([a-zA-Z0-9_+/#-]+)$`)

	// Ensure the topic length is within the allowed size
	if len(topic) > 65535 {
		return false
	}

	// Match the topic format
	return re.MatchString(topic)
}

func Encode(version ProtocolVersion, packet Packet) (data []byte, err error) {
	if version != 0 {
		packet.SetVersion(version)
	}
	return packet.Encode()
}

func Decode(version ProtocolVersion, data []byte) (packet Packet, err error) {
	var i int
	header := &FixedHeader{}
	n, err := header.Decode(data)
	if err != nil {
		return nil, err
	}
	i += n
	switch header.PacketType {
	case CONNECT:
		packet = NewConnect(header, version)
	case CONNACK:
		packet = NewConnack(header, version)
	case PUBLISH:
		packet = NewPublish(header, version)
	case PUBACK:
		packet = NewPuback(header, version)
	case PUBREC:
		packet = NewPubrec(header, version)
	case PUBREL:
		packet = NewPubrel(header, version)
	case PUBCOMP:
		packet = NewPubcomp(header, version)
	case SUBSCRIBE:
		packet = NewSubscribe(header, version)
	case SUBACK:
		packet = NewSuback(header, version)
	case UNSUBSCRIBE:
		packet = NewUnsubscribe(header, version)
	case UNSUBACK:
		packet = NewUnsuback(header, version)
	case PINGREQ:
		packet = NewPingreq(header, version)
	case PINGRESP:
		packet = NewPingresp(header, version)
	case DISCONNECT:
		packet = NewDisconnect(header, version)
	}
	if packet == nil {
		return nil, ErrMalformedPacket
	}
	if _, err = packet.Decode(data[i:]); err != nil {
		return
	}
	packet.SetVersion(version)
	return
}
