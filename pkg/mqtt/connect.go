package mqtt

func NewConnect(header *FixedHeader) *Connect {
	return &Connect{FixHeader: header}
}

type Connect struct {
	Version      ProtocolVersion
	FixHeader    *FixedHeader
	ProtocolName string
	//Clean Clean Session(v3,v4) or Clean Start(v5)
	UserNameFlag bool
	PasswordFlag bool
	WillFlag     bool
	Will         *Will
	Clean        bool
	KeepAlive    uint16
	ClientID     string
	Username     string
	Password     string
	Properties   *Properties
}

func (c *Connect) Encode() (result []byte, err error) {
	//Fixed header
	result = toHeader(CONNECT)
	//Variable header
	var buf []byte
	if buf, err = encodeUTF8Str(c.protocolName()); err != nil {
		return
	}
	result = append(result, buf...)
	result = append(result, byte(c.Version))
	result = append(result, c.encodeFlag())
	result = append(result, encodeUint16(c.KeepAlive)...)
	if buf, err = c.Properties.Encode(); err != nil {
		result = append(result, buf...)
	}
	if c.ClientID == "" {
		return nil, ErrNullClientID
	}
	if buf, err = encodeUTF8Str(c.ClientID); err != nil {
		return
	}
	result = append(result, buf...)
	if c.Will != nil {
		if c.Version == MQTT5 {
			// if buf, err = encodeWillProperties(c.Will.Properties); err != nil {
			// 	return
			// }
			result = append(result, buf...)
		}
		if buf, err = encodeUTF8Str(c.Will.Topic); err != nil {
			return
		}
		result = append(result, buf...)
		if buf, err = encodeUTF8Str(c.Will.Message); err != nil {
			return
		}
		result = append(result, buf...)
	}
	if buf, err = encodeUTF8Str(c.Username); err != nil {
		return
	}
	result = append(result, buf...)
	if buf, err = encodeUTF8Str(c.Password); err != nil {
		return
	}
	result = append(result, buf...)

	return
}

func (c *Connect) EncodeBody() ([]byte, error) {
	return nil, nil
}

func (c *Connect) Decode(data []byte) (int, error) {
	l, err := c.FixHeader.Decode(data)
	if err != nil {
		return 0, err
	}
	return c.DecodeBody(data[l:])
}

func (c *Connect) DecodeBody(data []byte) (int, error) {
	// total := len(data)
	var start int

	if protocol, n, err := decodeUTF8Str(data[start:]); err != nil {
		return start, err
	} else {
		c.ProtocolName = protocol
		start += n
	}

	if version, err := decodeProtocolVersion(data[start]); err != nil {
		return start, err
	} else {
		c.Version = version
		start++
	}
	{
		c.decodeFlag(data[start])
		start++
	}
	{
		c.KeepAlive = decodeKeepAlive(data[start : start+2])
		start += 2
	}

	if c.Version.IsMQTT5() {
		c.Properties = new(Properties)
		n, err := c.Properties.Decode(data[start:])
		if err != nil {
			return start, err
		}
		start += n
	}
	{
		clientID, n, err := decodeUTF8Str(data[start:])
		if err != nil {
			return start, err
		}
		c.ClientID = clientID
		start += n
	}
	if c.WillFlag {
		if c.Version.IsMQTT5() {
			c.Will.Properties.Decode(data[start:])
		}
		if topic, n, err := decodeUTF8Str(data[start:]); err != nil {
			return start, err
		} else {
			c.Will.Topic = topic
			start = start + n
		}
		if message, n, err := decodeUTF8Str(data[start:]); err != nil {
			return start, err
		} else {
			c.Will.Message = message
			start = start + n
		}
	}
	if c.UserNameFlag {
		if user, n, err := decodeUTF8Str(data[start:]); err != nil {
			return start, err
		} else {
			c.Username = user
			start = start + n
		}
	}
	if c.PasswordFlag {
		if pass, n, err := decodeUTF8Str(data[start:]); err != nil {
			return start, err
		} else {
			c.Password = pass
			start = start + n
		}
	}

	return start, nil
}

func (c *Connect) Read(r *Reader) (err error) {
	c.FixHeader = new(FixedHeader)
	if err := c.FixHeader.Read(r); err != nil {
		return err
	}
	return c.ReadBody(r)
}

func (c *Connect) ReadBody(r *Reader) error {
	data, err := r.Read(c.FixHeader.RemainingLength)
	if err != nil {
		return err
	}
	_, err = c.DecodeBody(data)
	return err
}

func (c *Connect) Write(reader *Writer) error {
	return nil
}

func (c *Connect) WriteBody(w *Writer) error {
	return nil
}

func (c *Connect) protocolName() string {
	switch c.Version {
	case MQTT3:
		return "MQIsdp"
	case MQTT4:
		return "MQTT"
	case MQTT5:
		return "MQTT"
	}
	return ""
}

func (c *Connect) decodeFlag(flag byte) {
	c.UserNameFlag = (flag&0b1000000 == 0b1000000)
	c.PasswordFlag = (flag&0b01000000 == 0b01000000)
	if flag&0b00000100 == 0b00000100 {
		c.Will = &Will{}
		c.Will.Retain = (flag&0b00100000 == 0b00100000)
		c.Will.QoS = QoS(flag & 0b00011000)
	}
	c.Clean = (flag&0b00000010 == 0b00000010)
	c.WillFlag = (flag&0b00000100 == 0b00000100)
	if c.WillFlag {
		c.Will = &Will{
			Retain:     flag&0b00100000 == 0b00100000,
			QoS:        QoS(flag & 0b00011000),
			Properties: new(Properties),
		}
	}
	return
}

func (c *Connect) encodeFlag() byte {
	var flags byte
	if c.Username != "" {
		flags = flags | 0b1
	}
	if c.Password != "" {
		flags = flags | 0b01
	}
	if c.Will != nil {
		flags = flags | (byte(c.Will.QoS) << 3)
		flags = flags | 0b00000010
	}
	if c.Clean {
		flags = flags | 0b00000001
	}
	return flags
}

func (c *Connect) String() string {
	return ""
}
