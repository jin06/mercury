package mqtt

type Connack struct {
	Version   ProtocolVersion
	FixHeader *FixedHeader
	// FixHeader *FixedHeader

	ReasonCode     ReasonCode
	Properties     *Properties
	SessionPresent bool
}

func (c *Connack) Encode() ([]byte, error) {
	// todo change size
	body, err := c.EncodeBody()
	if err != nil {
		return nil, err
	}
	c.FixHeader.RemainingLength = len(body)
	header, err := c.FixHeader.Encode()
	if err != nil {
		return nil, err
	}
	header = append(header, body...)
	return header, nil
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

func (c *Connack) Decode(data []byte) (int, error) {
	c.FixHeader = &FixedHeader{}
	n, err := c.FixHeader.Decode(data)
	if err != nil {
		return 0, err
	}
	bodyLen, err := c.DecodeBody(data[n:c.FixHeader.RemainingLength])
	return bodyLen + n, err
}

func (c *Connack) DecodeBody(data []byte) (int, error) {
	i, err := 0, error(nil)

	if c.SessionPresent, err = decodeBool(data[i]); err != nil {
		return i, err
	}
	i++
	c.ReasonCode = ReasonCode(data[i])
	i++
	if c.Properties == nil {
		c.Properties = &Properties{}
	}

	n, err := c.Properties.Decode(data[i:])
	if err != nil {
		return i, err
	}

	return i + n, nil
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
}

func (c *Connack) Write(w *Writer) error {
	data, err := c.Encode()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (c *Connack) WriteBody(w *Writer) error {
	data, err := c.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
