package mqtt

func NewDisconnect(header *FixedHeader, v ProtocolVersion) *Disconnect {
	return &Disconnect{BasePacket: &BasePacket{header, v}}
}

type Disconnect struct {
	*BasePacket
	ResionCode ReasonCode
	Properties *Properties
}

func (d *Disconnect) Encode() (result []byte, err error) {
	result = toHeader(DISCONNECT)
	if d.Version == MQTT5 {
		result = append(result, byte(d.ResionCode))
	}
	return
}

func (d *Disconnect) Decode(data []byte) (int, error) {
	if d.Version == MQTT5 {
		if len(data) < 1 {
			return 0, ErrBytesShorter
		}
		d.ResionCode = ReasonCode(data[0])
		return 1, nil
	}
	return len(data), nil
}

func (d *Disconnect) EncodeBody() ([]byte, error) {
	var data []byte
	if d.Version == MQTT5 {
		data = append(data, byte(d.ResionCode))
	}
	return data, nil
}

func (d *Disconnect) ReadBody(r *Reader) error {
	data, err := r.Read(d.Length())
	if err != nil {
		return err
	}
	_, err = d.Decode(data)
	return err
}

func (d *Disconnect) Write(w *Writer) error {
	data, err := d.Encode()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (d *Disconnect) WriteBody(w *Writer) error {
	data, err := d.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (d *Disconnect) PacketType() PacketType {
	return DISCONNECT
}

func (d *Disconnect) RemainingLength() int {
	if d.Version == MQTT5 {
		return 1
	}
	return 0
}

func (d *Disconnect) String() string {
	return "Disconnect Packet"
}

func (d *Disconnect) DecodeBody(data []byte) (int, error) {
	var start int

	// Decode Reason Code (MQTT 5.0 only)
	if d.Version == MQTT5 {
		if len(data) > start {
			d.ResionCode = ReasonCode(data[start])
			start++
		}

		// Decode Properties (MQTT 5.0 only)
		if len(data) > start {
			d.Properties = new(Properties)
			n, err := d.Properties.Decode(data[start:])
			if err != nil {
				return start, err
			}
			start += n
		}
	}

	return len(data), nil
}
