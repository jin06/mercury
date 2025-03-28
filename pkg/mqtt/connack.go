package mqtt

import "fmt"

func NewConnack(header *FixedHeader, v ProtocolVersion) *Connack {
	return &Connack{BasePacket: &BasePacket{FixedHeader: header, Version: v}}
}

type Connack struct {
	*BasePacket

	ReasonCode ReasonCode
	Properties *Properties
	// SessionPresent The Session Present flag informs the Client whether the Server is using Session State from a previous connection for this ClientID.
	SessionPresent bool
}

func (c *Connack) String() string {
	return fmt.Sprintf("Connack - ReasonCode: %d, Properties: %v, SessionPresent: %t", c.ReasonCode, c.Properties, c.SessionPresent)
}

func (c *Connack) Encode() ([]byte, error) {
	// todo change size
	body, err := c.EncodeBody()
	if err != nil {
		return nil, err
	}
	length := len(body)
	c.FixedHeader.RemainingLength = VariableByteInteger(length)
	header, err := c.FixedHeader.Encode()
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
	c.FixedHeader = &FixedHeader{}
	n, err := c.FixedHeader.Decode(data)
	if err != nil {
		return 0, err
	}
	bodyLen, err := c.DecodeBody(data[n:c.FixedHeader.RemainingLength])
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
	if err := c.BasePacket.FixedHeader.Read(r); err != nil {
		return err
	}
	return c.ReadBody(r)
}

func (c *Connack) ReadBody(r *Reader) error {
	data, err := r.Read(c.Length())
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
