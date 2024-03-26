package mqtt

import "io"

type Connack struct {
	ReasonCode ReasonCode
	Properties *ConnackProperties
}

func (c *Connack) Decode(reader io.Reader) (err error) {
	var msgLen byte
	if msgLen, err = readByte(reader); err != nil {
		return
	}
	length := int(msgLen)
	if length < 2 {
		return ProtocolError
	}
	if _, err = readByte(reader); err != nil {
		return
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
}

func decodeConnackProperties(reader io.Reader) (result *ConnackProperties, err error) {
	var total int
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
		case 0x27:
			{
				i += 4
				if result.MaximumPacketSize, err = readUint32(reader); err != nil {
					return
				}
			}
		case 0x25:
			{
				i++
				result.RetainAvailable, err = readBool(reader)
				if err != nil {
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
		case 0x29:
			{
				i++
				if result.SubscriptionIdentifierAvailable, err = readBool(reader); err != nil {
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
		case 0x22:
			{
				i += 2
				if result.TopicAliasMaximum, err = readUint16(reader); err != nil {
					return
				}
			}
		}
	}
	return
}
