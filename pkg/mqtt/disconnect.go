package mqtt

func NewDisconnect(header *FixedHeader) *Disconnect {
	return &Disconnect{FixedHeader: header}
}

type Disconnect struct {
	*FixedHeader
	Version    ProtocolVersion
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
	// ... implementation ...
	return nil
}

func (d *Disconnect) Write(w *Writer) error {
	// ... implementation ...
	return nil
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
	// Disconnect packet has no body, so just return 0 and nil
	return 0, nil
}
