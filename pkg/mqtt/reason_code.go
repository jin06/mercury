package mqtt

type ReasonCode byte

var (
	V3_CONNACK_ACCEPT                ReasonCode = 0x00
	V3_CONNACK_UNACCEPTABLE_VERSION  ReasonCode = 0x01
	V3_CONNACK_INDENTIFIER_REJECT    ReasonCode = 0x02
	V3_CONNACK_SERVER_UNAVAILABLE    ReasonCode = 0x03
	V3_CONNACK_BAD_USERNAME_PASSWORD ReasonCode = 0x04
	V3_CONNACK_NOT_AUTHORIZED        ReasonCode = 0x05
)

var (
	V5_CONACK_SUCCESS ReasonCode = 0x00
)
