package mqtt

type Version byte

func (v Version) String() string {
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
	MQTT3 Version = 0
	MQTT4 Version = 1
	MQTT5 Version = 2
)
