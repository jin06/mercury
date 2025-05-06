package mqtt

import "slices"

const (
	MQTT3 ProtocolVersion = 3
	MQTT4 ProtocolVersion = 4
	MQTT5 ProtocolVersion = 5

	PayloadFormatBytes  PayloadFormat = 0
	PayloadFormatString PayloadFormat = 1
)

var supportedVersions = []ProtocolVersion{
	MQTT3, MQTT4, MQTT5,
}

type PayloadFormat = byte

type ProtocolVersion byte

func (v ProtocolVersion) String() string {
	switch v {
	case MQTT3:
		return "3"
	case MQTT4:
		return "4"
	case MQTT5:
		return "5"
	}
	return "unsupported"
}

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

func (v ProtocolVersion) Valid() error {
	// switch v {
	// case MQTT3:
	// case MQTT4:
	// case MQTT5:
	// default:
	// 	return ErrUnsupportVersion
	// }

	if !slices.Contains(supportedVersions, v) {
		return ErrUnsupportVersion
	}

	return nil
}

func encodeProtocolVersion(v ProtocolVersion) (byte, error) {
	if err := v.Valid(); err != nil {
		return 0, err
	}
	return byte(v), nil
}

func decodeProtocolVersion(b byte) (ProtocolVersion, error) {
	// switch b {
	// case 3:
	// 	return MQTT3, nil
	// case 4:
	// 	return MQTT4, nil
	// case 5:
	// 	return MQTT5, nil
	// }
	// return 0, ErrUnsupportVersion
	v := ProtocolVersion(b)
	err := v.Valid()
	return v, err
}
