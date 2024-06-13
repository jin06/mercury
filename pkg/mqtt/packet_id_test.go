package mqtt

import (
	"testing"
)

func TestPacketID(t *testing.T) {
	var pid PacketID = 128 + 256
	bytes := pid.ToBytes()
	t.Log(bytes)
	if bytes[0] != 0b00000001 {
		t.Fail()
	}
	if bytes[1] != 0b10000000 {
		t.Fail()
	}
}
