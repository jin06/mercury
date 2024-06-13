package mqtt

import "io"

type Pingreq struct {
}

func (p *Pingreq) Encode() (result []byte, err error) {
	result = toHeader(PINGREQ)
	return
}

func (p *Pingreq) Decode(reader io.Reader) error {
	return nil
}
