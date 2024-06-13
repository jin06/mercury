package mqtt

import "io"

type Pingresp struct {
}

func (p *Pingresp) Encode() ([]byte, error) {
	return toHeader(PINGRESP), nil
}

func (p *Pingresp) Decode(reader io.Reader) error {
	return nil
}
