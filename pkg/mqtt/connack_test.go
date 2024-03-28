package mqtt

import "testing"

func TestEncodeConnack(t *testing.T) {
	v5data := []byte{
		0x20, 0x16, 0x00, 0x00, 0x13, 0x27, 0x00, 0x10, 0x00, 0x00, 0x21, 0x00, 0x20, 0x25, 0x01, 0x2a, 0x01, 0x29, 0x01, 0x22, 0xff, 0xff, 0x28, 0x01,
	}
	t.Log(v5data)
	t.Log(len(v5data))
	ack := Connack{
		Properties: &ConnackProperties{
			MaximumPacketSize:               1048576,
			TopicAliasMaximum:               65535,
			RetainAvailable:                 true,
			SharedSubscriptionAvailable:     true,
			SubscriptionIdentifierAvailable: true,
			WildcardSubscriptionAvailable:   true,
			SessionExpiryInterval:           65535,
			ReceiveMaximum:                  32,
		},
	}
	v5result, err := ack.Encode(MQTT5)
	if err != nil {
		t.Error(err)
	}
	t.Log(v5result)
	t.Log(len(v5result))
	if len(v5data) != len(v5result) {
		t.Fail()
	}
	for i := 0; i < len(v5data); i++ {

		if v5data[i] != v5result[i] {
			t.Fail()
		}
	}
}
