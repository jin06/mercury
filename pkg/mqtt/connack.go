package mqtt

type Connack struct {
	Version   ProtocolVersion
	FixHeader *FixedHeader
	// FixHeader *FixedHeader

	ReasonCode     ReasonCode
	Properties     *Properties
	SessionPresent bool
}

func (c *Connack) Encode() (result []byte, err error) {
	// todo change size
	if result, err = c.FixHeader.Encode(); err != nil {
		return nil, err
	}
	if body, err := c.EncodeBody(); err != nil {
		return nil, err
	} else {
		result = append(result, body...)
	}
	return
}

func (c *Connack) EncodeBody() ([]byte, error) {
	result := make([]byte, 0)
	result = append(result, encodeBool(c.SessionPresent))
	result = append(result, byte(c.ReasonCode))
	if data, err := c.Properties.Encode(); err != nil {
		return nil, err
	} else {
		result = append(result, data...)
	}
	return result, nil
}

func (c *Connack) Decode(data []byte) (n int, err error) {
	return
}

func (c *Connack) DecodeBody(data []byte) (n int, err error) {
	return
}

func (c *Connack) Read(r *Reader) error {
	c.FixHeader = new(FixedHeader)
	if err := c.FixHeader.Read(r); err != nil {
		return err
	}
	return c.ReadBody(r)
}

func (c *Connack) ReadBody(r *Reader) error {
	data, err := r.Read(c.FixHeader.RemainingLength)
	if err != nil {
		return err
	}
	_, err = c.DecodeBody(data)
	return err
	// var msgLen byte
	// if msgLen, err = readByte(reader); err != nil {
	// 	return
	// }
	// length := int(msgLen)
	// if length < 2 {
	// 	return ErrProtocol
	// }
	// if flags, err := readByte(reader); err != nil {
	// 	if (flags & 0x00000001) == 0x00000001 {
	// 		c.SessionPresent = true
	// 	}
	// 	return err
	// }
	// if code, err := readByte(reader); err != nil {
	// 	return err
	// } else {
	// 	c.ReasonCode = ReasonCode(code)
	// }
	// if length > 2 {
	// 	c.Properties, err = readProperties(reader)
	// }
}

func (c *Connack) Write(w *Writer) error {
	_, err := w.Write([]byte{encodeBool(c.SessionPresent)})
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
