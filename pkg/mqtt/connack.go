package mqtt

type Connack struct {
	Version   ProtocolVersion
	FixHeader *FixedHeader

	ReasonCode     ReasonCode
	Properties     *Properties
	SessionPresent bool
}

func (c *Connack) Encode() (result []byte, err error) {
	// result = toHeader(CONNACK)
	// result = append(result, 0, byte(c.ReasonCode))
	// if c.SessionPresent {
	// 	result[2] |= 0x00000001
	// }
	// if c.Version == MQTT5 {
	// 	result = append(result, 24)
	// 	var pl int = 0
	// 	result = append(result, 0x11)
	// 	pl++
	// 	result = append(result, uint32ToBytes(c.Properties.SessionExpiryInterval)...)
	// 	pl += 4
	// 	result = append(result, 0x21)
	// 	pl++
	// 	result = append(result, uint16ToBytes(c.Properties.ReceiveMaximum)...)
	// 	pl += 2
	// 	result = append(result, 0x22)
	// 	pl++
	// 	result = append(result, uint16ToBytes(c.Properties.TopicAliasMaximum)...)
	// 	pl += 2
	// 	result = append(result, 0x25)
	// 	pl++
	// 	result = append(result, boolTobyte(c.Properties.RetainAvailable))
	// 	pl++
	// 	result = append(result, 0x2a)
	// 	pl++
	// 	result = append(result, boolTobyte(c.Properties.SharedSubscriptionAvailable))
	// 	pl++
	// 	result = append(result, 0x27)
	// 	pl++
	// 	result = append(result, uint32ToBytes(c.Properties.MaximumPacketSize)...)
	// 	pl += 4
	// 	result = append(result, 0x28)
	// 	pl++
	// 	result = append(result, boolTobyte(c.Properties.WildcardSubscriptionAvailable))
	// 	pl++
	// 	result = append(result, 0x29)
	// 	pl++
	// 	result = append(result, boolTobyte(c.Properties.SubscriptionIdentifierAvailable))
	// 	pl++
	// 	result[4] = byte(pl)
	// 	result[1] = byte(pl + 3)
	// }

	return
}

func (c *Connack) EncodeBody() ([]byte, error) {
	return nil, nil
}

func (c *Connack) Decode(data []byte) (n int, err error) {
	return
}

func (c *Connack) DecodeBody(data []byte) (n int, err error) {
	return
}

func (c *Connack) Read(reader *Reader) (err error) {
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
		c.Properties, err = readProperties(reader)
	}
	return
}

func (c *Connack) ReadBody(r *Reader) (err error) {
	return
}

func (c *Connack) Write(w *Writer) error {
	_, err := w.Write([]byte{boolTobyte(c.SessionPresent)})
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{byte(c.ReasonCode)})
	if err != nil {
		return err
	}
	if c.Properties != nil {
		writeProperties(w, c.Properties)
	}
	return nil
}

func (c *Connack) WriteBody(w *Writer) error {
	return nil
}
