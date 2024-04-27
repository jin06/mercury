package mqtt

type ProtocolVersion byte

func (v ProtocolVersion) String() string {
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

const (
	MinMessageID uint16 = 0
	MaxMessageID uint16 = 65535
)

type TopicWildcard string
