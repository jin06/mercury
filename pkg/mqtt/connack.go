package mqtt

import (
	"io"
)

type Connack struct {
	ReasonCode     ReasonCode
	Properties     *ConnackProperties
	SessionPresent bool
}

func (c *Connack) Decode(reader io.Reader) (err error) {
	var msgLen byte
	if msgLen, err = readByte(reader); err != nil {
		return
	}
	length := int(msgLen)
	if length < 2 {
		return ErrProtocol
	}
	if flags, err := readByte(reader); err != nil {
		if (flags & 0x00000001) == 0x00000001 {
			c.SessionPresent = true
		}
		return err
	}
	if code, err := readByte(reader); err != nil {
		return err
	} else {
		c.ReasonCode = ReasonCode(code)
	}
	if length > 2 {
		c.Properties, err = decodeConnackProperties(reader)
	}
	return
}

type ConnackProperties struct {
	MaximumPacketSize               uint32
	RetainAvailable                 bool
	SharedSubscriptionAvailable     bool
	SubscriptionIdentifierAvailable bool
	TopicAliasMaximum               uint16
	WildcardSubscriptionAvailable   bool
	ReceiveMaximum                  uint16
	SessionExpiryInterval           uint32
}

func decodeConnackProperties(reader io.Reader) (result *ConnackProperties, err error) {
	var total int
	result = &ConnackProperties{}
	if res, err := readByte(reader); err != nil {
		return result, err
	} else {
		total = int(res)
	}

	for i := 0; i < total; {
		var identifier byte
		if identifier, err = readByte(reader); err != nil {
			return
		}
		i++

		switch identifier {
		case 0x11:
			{
				i += 4
				if result.SessionExpiryInterval, err = readUint32(reader); err != nil {
					return
				}
			}
		case 0x21:
			{
				i += 2
				if result.ReceiveMaximum, err = readUint16(reader); err != nil {
					return
				}
			}
		case 0x22:
			{
				i += 2
				if result.TopicAliasMaximum, err = readUint16(reader); err != nil {
					return
				}
			}
		case 0x25:
			{
				i++
				if result.RetainAvailable, err = readBool(reader); err != nil {
					return
				}
			}
		case 0x27:
			{
				i += 4
				if result.MaximumPacketSize, err = readUint32(reader); err != nil {
					return
				}
			}
		case 0x28:
			{
				i++
				if result.WildcardSubscriptionAvailable, err = readBool(reader); err != nil {
					return
				}
			}
		case 0x29:
			{
				i++
				if result.SubscriptionIdentifierAvailable, err = readBool(reader); err != nil {
					return
				}
			}
		case 0x2a:
			{
				i++
				result.SharedSubscriptionAvailable, err = readBool(reader)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func (c *Connack) Encode(version ProtocolVersion) (result []byte, err error) {
	result = toHeader(CONNACK)
	result = append(result, 0, byte(c.ReasonCode))
	if c.SessionPresent {
		result[2] |= 0x00000001
	}
	if version == MQTT5 {
		result = append(result, 24)
		var pl int = 0
		result = append(result, 0x11)
		pl++
		result = append(result, uint32ToBytes(c.Properties.SessionExpiryInterval)...)
		pl += 4
		result = append(result, 0x21)
		pl++
		result = append(result, uint16ToBytes(c.Properties.ReceiveMaximum)...)
		pl += 2
		result = append(result, 0x22)
		pl++
		result = append(result, uint16ToBytes(c.Properties.TopicAliasMaximum)...)
		pl += 2
		result = append(result, 0x25)
		pl++
		result = append(result, boolTobyte(c.Properties.RetainAvailable))
		pl++
		result = append(result, 0x2a)
		pl++
		result = append(result, boolTobyte(c.Properties.SharedSubscriptionAvailable))
		pl++
		result = append(result, 0x27)
		pl++
		result = append(result, uint32ToBytes(c.Properties.MaximumPacketSize)...)
		pl += 4
		result = append(result, 0x28)
		pl++
		result = append(result, boolTobyte(c.Properties.WildcardSubscriptionAvailable))
		pl++
		result = append(result, 0x29)
		pl++
		result = append(result, boolTobyte(c.Properties.SubscriptionIdentifierAvailable))
		pl++
		result[4] = byte(pl)
		result[1] = byte(pl + 3)
	}

	return
}
