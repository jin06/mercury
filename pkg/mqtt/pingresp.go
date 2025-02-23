package mqtt

func NewPingresp(header *FixedHeader) *Pingresp {
	return &Pingresp{FixedHeader: header}
}

type Pingresp struct {
	*FixedHeader
}

func (p *Pingresp) Encode() ([]byte, error) {
	p.FixedHeader.RemainingLength = 0
	return p.FixedHeader.Encode()
}

func (p *Pingresp) Decode(data []byte) (int, error) {
	// No body to decode for Pingresp
	return len(data), nil
}

func (p *Pingresp) DecodeBody(data []byte) (int, error) {
	// No body to decode for Pingresp
	return len(data), nil
}

func (p *Pingresp) ReadBody(r *Reader) error {
	return nil
}

func (p *Pingresp) Write(w *Writer) error {
	data, err := p.Encode()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (p *Pingresp) EncodeBody() ([]byte, error) {
	// No body to encode for Pingresp
	return []byte{}, nil
}

func (p *Pingresp) WriteBody(w *Writer) error {
	data, err := p.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (p *Pingresp) PacketType() PacketType {
	return PINGRESP
}

func (p *Pingresp) RemainingLength() int {
	return 0
}

func (p *Pingresp) String() string {
	return "Pingresp Packet"
}
