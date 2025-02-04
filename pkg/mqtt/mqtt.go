package mqtt

type ProtocolVersion byte

func (v ProtocolVersion) IsMQTT5() bool {
	return v == MQTT5
}

func (v ProtocolVersion) Name() string {
	switch v {
	case MQTT3:
		return "mqtt3.1"
	case MQTT4:
		return "mqtt3.1.1"
	case MQTT5:
		return "mqtt5"
	}
	return "unsupported"
}

const (
	MQTT3 ProtocolVersion = 3
	MQTT4 ProtocolVersion = 4
	MQTT5 ProtocolVersion = 5
)

type TopicWildcard string
