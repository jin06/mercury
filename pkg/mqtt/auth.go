package mqtt

func NewAuth(header *FixedHeader, v ProtocolVersion) *Auth {
	// return &Auth{FixedHeader: header}
	return &Auth{
		BasePacket: &BasePacket{header, v},
		Properties: new(Properties),
	}
}

type Auth struct {
	*BasePacket
	// Version    ProtocolVersion
	ReasonCode ReasonCode
	Properties *Properties
}

func (a *Auth) Encode() ([]byte, error) {
	body, err := a.EncodeBody()
	if err != nil {
		return nil, err
	}
	a.FixedHeader.RemainingLength = VariableByteInteger(len(body))
	header, err := a.FixedHeader.Encode()
	if err != nil {
		return nil, err
	}
	return append(header, body...), nil
}

func (a *Auth) Decode(data []byte) (int, error) {
	n, err := a.FixedHeader.Decode(data)
	if err != nil {
		return 0, err
	}
	bodyLen, err := a.DecodeBody(data[n:])
	return bodyLen + n, err
}

func (a *Auth) ReadBody(r *Reader) error {
	data, err := r.Read(a.Length())
	if err != nil {
		return err
	}
	_, err = a.DecodeBody(data)
	return err
}

func (a *Auth) Write(w *Writer) error {
	data, err := a.Encode()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (a *Auth) PacketType() PacketType {
	return AUTH
}

func (a *Auth) RemainingLength() int {
	length := 1 // Reason Code length
	if a.Version == MQTT5 && a.Properties != nil {
		propertiesLength, _ := a.Properties.Encode()
		length += len(propertiesLength)
	}
	return length
}

func (a *Auth) String() string {
	return "Auth Packet"
}

func (a *Auth) DecodeBody(data []byte) (int, error) {
	var start int

	// Decode Reason Code (MQTT 5.0 only)
	if a.Version == MQTT5 {
		if len(data) > start {
			a.ReasonCode = ReasonCode(data[start])
			start++
		}

		// Decode Properties (MQTT 5.0 only)
		if len(data) > start {
			a.Properties = new(Properties)
			n, err := a.Properties.Decode(data[start:])
			if err != nil {
				return start, err
			}
			start += n
		}
	}

	return len(data), nil
}

func (a *Auth) EncodeBody() ([]byte, error) {
	var data []byte

	// Encode Reason Code (MQTT 5.0 only)
	if a.Version == MQTT5 {
		data = append(data, byte(a.ReasonCode))

		// Encode Properties (MQTT 5.0 only)
		if a.Properties != nil {
			propertiesData, err := a.Properties.Encode()
			if err != nil {
				return nil, err
			}
			data = append(data, propertiesData...)
		}
	}

	return data, nil
}

func (a *Auth) WriteBody(w *Writer) error {
	data, err := a.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
