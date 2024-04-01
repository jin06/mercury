package mqtt

import (
	"bytes"
	"testing"
)

func TestEncodeConnack(t *testing.T) {
	ack := Connack{
		ReasonCode: V5_CONACK_SUCCESS,
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
		SessionPresent: true,
	}
	var err error
	var encode []byte
	if encode, err = ack.Encode(MQTT5); err != nil {
		t.Error(err)
	}
	reader := bytes.NewReader(encode)
	reader.Read(make([]byte, 1))
	ack2 := Connack{
		Properties: &ConnackProperties{},
	}
	if err = ack2.Decode(reader); err != nil {
		t.Error(err)
	}
	if ack.ReasonCode != ack2.ReasonCode {
		t.Fail()
	}
	if ack.SessionPresent != ack2.SessionPresent {
		t.Fail()
	}
}
