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
	MQTT3 ProtocolVersion = 0
	MQTT4 ProtocolVersion = 1
	MQTT5 ProtocolVersion = 2
)
