package mqtt

func NewPingreq(header *FixedHeader) *Pingreq {
	return &Pingreq{FixedHeader: header}
}

type Pingreq struct {
	*FixedHeader
}

func (p *Pingreq) Response() *Pingresp {
	return &Pingresp{
		FixedHeader: &FixedHeader{
			PacketType:      PINGRESP,
			Flags:           0,
			RemainingLength: 0,
		},
	}
}

func (p *Pingreq) Encode() (result []byte, err error) {
	result = toHeader(PINGREQ)
	return
}

func (p *Pingreq) Decode(data []byte) (int, error) {
	// No body to decode for Pingreq
	return len(data), nil
}

func (p *Pingreq) DecodeBody(data []byte) (int, error) {
	// No body to decode for Pingreq
	return len(data), nil
}

func (p *Pingreq) ReadBody(r *Reader) error {
	return nil
}

func (p *Pingreq) Write(w *Writer) error {
	return nil
}

func (p *Pingreq) PacketType() PacketType {
	return PINGREQ
}

func (p *Pingreq) RemainingLength() int {
	return 0
}

func (p *Pingreq) String() string {
	return "Pingreq Packet"
}

func (p *Pingreq) Read(r *Reader) error {
	return p.ReadBody(r)
}

func (p *Pingreq) EncodeBody() ([]byte, error) {
	// No body to encode for Pingreq
	return []byte{}, nil
}

func (p *Pingreq) WriteBody(w *Writer) error {
	data, err := p.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
